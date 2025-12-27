package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"shopping/internal/web/views"
)

func (s *Server) handleShoppingListPage(w http.ResponseWriter, r *http.Request) {
	user, _ := s.auth.CurrentUser(r)

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	items, err := s.shopping.svc.ListItems(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	data := views.ShoppingListPageData{
		Base: views.BaseData{
			Title:         "Lista zakupów",
			User:          user,
			HTMXSrc:       s.cfg.HTMXSrc,
			StaticVersion: s.staticV,
			IsAdmin:       s.isAdmin(user),
		},
		Items: items,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.ShoppingListPage(data).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleShoppingListPartial(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	items, err := s.shopping.svc.ListItems(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.ShoppingListCard(views.ShoppingListData{Items: items}).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleAddShoppingListByName(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	qty, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("quantity")), 64)
	if err != nil {
		http.Error(w, "Nieprawidłowa ilość.", http.StatusBadRequest)
		return
	}
	unit := products.Unit(strings.TrimSpace(r.FormValue("unit")))

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if err := s.shopping.svc.AddItemByName(ctx, name, qty, unit); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.events.Publish(eventShoppingList, clientIDFromRequest(r))
	s.shoppingListResponse(w, r)
}

func (s *Server) handleAddShoppingListFromProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathInt64(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if err := s.shopping.svc.AddItemByProductID(ctx, id); err != nil {
		s.writeDBError(w, err)
		return
	}
	s.events.Publish(eventShoppingList, clientIDFromRequest(r))
	if r.Header.Get("HX-Request") == "true" {
		if strings.Contains(r.Header.Get("HX-Current-URL"), "/shopping-list") {
			s.handleShoppingListPartial(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if ref := strings.TrimSpace(r.Referer()); ref != "" {
		http.Redirect(w, r, ref, http.StatusFound)
		return
	}
	http.Redirect(w, r, "/products", http.StatusFound)
}

func (s *Server) handleSetShoppingListDone(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathInt64(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	doneStr := strings.TrimSpace(r.FormValue("done"))
	qtyStr := strings.TrimSpace(r.FormValue("quantity"))
	unit := products.Unit(strings.TrimSpace(r.FormValue("unit")))

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if qtyStr != "" || unit != "" {
		qty, err := strconv.ParseFloat(qtyStr, 64)
		if err != nil {
			http.Error(w, "Nieprawidłowa ilość.", http.StatusBadRequest)
			return
		}
		if unit == "" {
			unit = products.UnitPiece
		}
		if err := s.shopping.svc.SetQuantity(ctx, shoppinglist.ItemID(id), qty, unit); err != nil {
			s.writeUserError(w, err)
			return
		}
	}
	if doneStr != "" {
		done := doneStr == "1" || strings.EqualFold(doneStr, "true") || strings.EqualFold(doneStr, "on")
		if err := s.shopping.svc.SetDone(ctx, shoppinglist.ItemID(id), done); err != nil {
			s.writeDBError(w, err)
			return
		}
		if done {
			s.events.Publish(eventProductsList, clientIDFromRequest(r))
		}
	}
	s.events.Publish(eventShoppingList, clientIDFromRequest(r))
	s.shoppingListResponse(w, r)
}

func (s *Server) handleDeleteShoppingListItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathInt64(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if err := s.shopping.svc.Delete(ctx, shoppinglist.ItemID(id)); err != nil {
		s.writeDBError(w, err)
		return
	}
	s.events.Publish(eventShoppingList, clientIDFromRequest(r))
	s.shoppingListResponse(w, r)
}

func (s *Server) handleAddShoppingItemToSupplies(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathInt64(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	item, err := s.shopping.svc.GetItem(ctx, shoppinglist.ItemID(id))
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	if item.ProductID == nil {
		if existingID, found, err := s.shopping.svc.FindProductIDByName(ctx, item.Name); err != nil {
			s.writeDBError(w, err)
			return
		} else if found {
			if err := s.shopping.svc.LinkToProduct(ctx, shoppinglist.ItemID(id), existingID, item.Name); err != nil {
				s.writeDBError(w, err)
				return
			}
		} else {
			pid, err := s.products.svc.CreateProduct(ctx, products.NewProduct{
				Name:     item.Name,
				Quantity: item.Quantity,
				Unit:     item.Unit,
			})
			if err != nil {
				s.writeUserError(w, err)
				return
			}
			if err := s.shopping.svc.LinkToProduct(ctx, shoppinglist.ItemID(id), int64(pid), item.Name); err != nil {
				s.writeDBError(w, err)
				return
			}
		}
	}

	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	s.events.Publish(eventShoppingList, clientIDFromRequest(r))
	s.shoppingListResponse(w, r)
}

func parsePathInt64(r *http.Request, key string) (int64, bool) {
	raw := strings.TrimSpace(r.PathValue(key))
	if raw == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func (s *Server) shoppingListResponse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		s.handleShoppingListPartial(w, r)
		return
	}
	http.Redirect(w, r, "/shopping-list", http.StatusFound)
}
