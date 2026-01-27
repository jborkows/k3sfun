package web

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"shopping/internal/domain/products"
	"shopping/internal/web/views"
)

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	topic := eventTopic(strings.TrimSpace(r.URL.Query().Get("topic")))
	if topic != eventShoppingList && topic != eventProductsList {
		http.Error(w, "bad topic", http.StatusBadRequest)
		return
	}
	clientID := strings.TrimSpace(r.URL.Query().Get("client"))

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	sub, unsubscribe := s.events.Subscribe(topic, clientID)
	defer unsubscribe()

	// No initial update - page already has data rendered server-side.
	// Only send updates when events occur.

	keepalive := time.NewTicker(20 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-keepalive.C:
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case <-sub:
			if err := s.writeEventUpdate(r.Context(), w, topic, r); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func (s *Server) writeEventUpdate(ctx context.Context, w http.ResponseWriter, topic eventTopic, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var payload []byte
	var err error
	switch topic {
	case eventShoppingList:
		editMode, shortMode := parseShoppingView(r)
		if !editMode && !shortMode {
			shortMode = true
		}
		payload, err = s.renderShoppingListHTML(ctx, editMode, shortMode)
	case eventProductsList:
		payload, err = s.renderProductsListHTML(ctx, r)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
	if err != nil {
		return err
	}

	return writeSSE(w, string(topic), payload)
}

func (s *Server) renderShoppingListHTML(ctx context.Context, editMode bool, shortMode bool) ([]byte, error) {
	items, err := s.shopping.svc.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := views.ShoppingListCard(views.ShoppingListData{Items: items, Units: s.units, EditMode: editMode, ShortMode: shortMode}).Render(ctx, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Server) renderProductsListHTML(ctx context.Context, r *http.Request) ([]byte, error) {
	source := r
	q := r.URL.Query()
	hasFilters := q.Has("missing") || q.Has("q") || q.Has("groups") || q.Has("page")
	if !hasFilters {
		if ref := strings.TrimSpace(r.Referer()); ref != "" {
			if u, err := url.Parse(ref); err == nil {
				r2 := r.Clone(r.Context())
				r2.URL = u
				source = r2
			}
		}
	}

	onlyMissing, nameQuery, groupNames, formGroupIDs, page := parseProductsListQuery(source)
	perPage := products.MaxProductsPageSize

	groups, err := s.products.qry.ListGroups(ctx)
	if err != nil {
		return nil, err
	}

	// Resolve group names to IDs (from URL), and merge with form-submitted IDs
	groupIDs := resolveGroupNames(groupNames, groups)
	if len(groupIDs) == 0 {
		groupIDs = formGroupIDs
	}

	total, err := s.products.qry.CountProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: onlyMissing,
		NameQuery:        nameQuery,
		GroupIDs:         groupIDs,
	})
	if err != nil {
		return nil, err
	}

	totalPages := int64(1)
	if total > 0 {
		totalPages = (total + perPage - 1) / perPage
	}
	if page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * perPage
	list, err := s.products.qry.ListProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: onlyMissing,
		NameQuery:        nameQuery,
		GroupIDs:         groupIDs,
		Limit:            perPage,
		Offset:           offset,
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := views.ProductsList(views.ProductsListData{
		Groups:           groups,
		Products:         list,
		OnlyMissing:      onlyMissing,
		NameQuery:        nameQuery,
		SelectedGroupIDs: groupIDs,
		Page:             page,
		TotalPages:       totalPages,
		Total:            total,
	}).Render(ctx, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSSE(w http.ResponseWriter, event string, data []byte) error {
	if _, err := fmt.Fprintf(w, "event: %s\n", event); err != nil {
		return err
	}
	for line := range bytes.SplitSeq(data, []byte("\n")) {
		if _, err := w.Write([]byte("data: ")); err != nil {
			return err
		}
		if _, err := w.Write(line); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\n"))
	return err
}
