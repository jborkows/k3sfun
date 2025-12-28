package shoppinglist

import (
	"context"
	"log/slog"
	"time"

	"shopping/internal/domain/products"
)

type Service struct {
	repo        Repository
	productsSvc *products.Service
}

func NewService(repo Repository, productsSvc *products.Service) *Service {
	return &Service{repo: repo, productsSvc: productsSvc}
}

func (s *Service) ListItems(ctx context.Context) ([]Item, error) {
	if err := s.repo.CleanupDoneBefore(ctx, time.Now().Add(-6*time.Hour)); err != nil {
		return nil, err
	}
	return s.repo.ListItems(ctx)
}

func (s *Service) GetItem(ctx context.Context, id ItemID) (Item, error) {
	return s.repo.GetItem(ctx, id)
}

func (s *Service) AddItemByName(ctx context.Context, name string, qty products.Quantity, unit products.Unit) error {
	normalized, err := NormalizeItemName(name)
	if err != nil {
		return err
	}
	if qty <= 0 {
		return ErrQuantityMustBePositive
	}
	if unit == "" {
		unit = products.UnitPiece
	}
	if _, err := products.NormalizeUnit(unit); err != nil {
		return err
	}
	return s.repo.AddItemByName(ctx, normalized, qty, unit)
}

func (s *Service) AddItemByProductID(ctx context.Context, productID int64) error {
	return s.repo.AddItemByProductID(ctx, productID)
}

func (s *Service) SetDone(ctx context.Context, id ItemID, done bool) error {
	if !done || s.productsSvc == nil {
		return s.repo.SetDone(ctx, id, done)
	}

	item, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return err
	}
	if item.Done == done {
		return nil
	}

	qty := item.Quantity
	if qty < 0 {
		qty = 0
	}

	if item.ProductID == nil {
		if existingID, found, err := s.repo.FindProductIDByName(ctx, item.Name); err != nil {
			return err
		} else if found {
			if err := s.repo.LinkToProduct(ctx, id, existingID, item.Name); err != nil {
				return err
			}
			pid := products.ProductID(existingID)
			item.ProductID = &pid
		} else {
			pid, err := s.productsSvc.CreateProduct(ctx, products.NewProduct{
				Name:     item.Name,
				Quantity: 0,
				Unit:     item.Unit,
			})
			if err != nil {
				return err
			}
			if err := s.repo.LinkToProduct(ctx, id, int64(pid), item.Name); err != nil {
				return err
			}
			item.ProductID = &pid
		}
	}

	if item.ProductID != nil && qty > 0 {
		if err := s.productsSvc.AddProductQuantity(ctx, *item.ProductID, qty); err != nil {
			return err
		}
		if err := s.productsSvc.SetProductMissing(ctx, *item.ProductID, false); err != nil {
			slog.Warn("failed to clear missing flag after adding quantity",
				"product_id", *item.ProductID,
				"error", err,
			)
		}
	}
	return s.repo.SetDone(ctx, id, done)
}

func (s *Service) SetQuantity(ctx context.Context, id ItemID, qty products.Quantity, unit products.Unit) error {
	if qty <= 0 {
		return ErrQuantityMustBePositive
	}
	if _, err := products.NormalizeUnit(unit); err != nil {
		return err
	}
	return s.repo.SetQuantity(ctx, id, qty, unit)
}

func (s *Service) Delete(ctx context.Context, id ItemID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) LinkToProduct(ctx context.Context, id ItemID, productID int64, name string) error {
	return s.repo.LinkToProduct(ctx, id, productID, name)
}

func (s *Service) FindProductIDByName(ctx context.Context, name string) (int64, bool, error) {
	return s.repo.FindProductIDByName(ctx, name)
}
