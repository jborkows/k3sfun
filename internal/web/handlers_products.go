package web

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"shopping/internal/domain/products"
	"shopping/internal/web/views"
)

// productsListData contains the shared data needed for both full page and partial renders.
type productsListData struct {
	Groups           []products.Group
	Products         []products.Product
	OnlyMissing      bool
	NameQuery        string
	SelectedGroupIDs []products.GroupID
	Page             int64
	TotalPages       int64
	Total            int64
}

// fetchProductsListData fetches all data needed to render the products list.
// This consolidates the duplicated logic from handleProductsPage and handleProductsPartial.
func (s *Server) fetchProductsListData(ctx context.Context, r *http.Request) (*productsListData, error) {
	onlyMissing, nameQuery, groupNames, formGroupIDs, page := parseProductsListQuery(r)
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

	filter := products.ProductFilter{
		OnlyMissingOrLow: onlyMissing,
		NameQuery:        nameQuery,
		GroupIDs:         groupIDs,
	}

	total, err := s.products.qry.CountProducts(ctx, filter)
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
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * perPage
	filter.Limit = perPage
	filter.Offset = offset

	list, err := s.products.qry.ListProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &productsListData{
		Groups:           groups,
		Products:         list,
		OnlyMissing:      onlyMissing,
		NameQuery:        nameQuery,
		SelectedGroupIDs: groupIDs,
		Page:             page,
		TotalPages:       totalPages,
		Total:            total,
	}, nil
}

func (s *Server) handleProductsNewPage(w http.ResponseWriter, r *http.Request) {
	user, _ := s.auth.CurrentUser(r)

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	groups, err := s.products.qry.ListGroups(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	pageData := views.ProductsNewPageData{
		Base: views.BaseData{
			Title:         "Dodaj produkt / grupę",
			User:          user,
			HTMXSrc:       s.cfg.HTMXSrc,
			StaticVersion: s.staticV,
			IsAdmin:       s.isAdmin(user),
		},
		Groups: groups,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.ProductsNewPage(pageData).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleProductsPage(w http.ResponseWriter, r *http.Request) {
	user, _ := s.auth.CurrentUser(r)

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	data, err := s.fetchProductsListData(ctx, r)
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	pageData := views.ProductsPageData{
		Base: views.BaseData{
			Title:         "Zapasy",
			User:          user,
			HTMXSrc:       s.cfg.HTMXSrc,
			StaticVersion: s.staticV,
			IsAdmin:       s.isAdmin(user),
		},
		Groups:           data.Groups,
		Products:         data.Products,
		OnlyMissing:      data.OnlyMissing,
		NameQuery:        data.NameQuery,
		SelectedGroupIDs: data.SelectedGroupIDs,
		Page:             data.Page,
		TotalPages:       data.TotalPages,
		Total:            data.Total,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.ProductsPage(pageData).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleProductsPartial(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	data, err := s.fetchProductsListData(ctx, r)
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	listData := views.ProductsListData{
		Groups:           data.Groups,
		Products:         data.Products,
		OnlyMissing:      data.OnlyMissing,
		NameQuery:        data.NameQuery,
		SelectedGroupIDs: data.SelectedGroupIDs,
		Page:             data.Page,
		TotalPages:       data.TotalPages,
		Total:            data.Total,
	}

	// Set HX-Push-Url header so browser URL updates to match filter state
	pushURL := buildProductsPageURL(data.OnlyMissing, data.NameQuery, data.SelectedGroupIDs, data.Groups, data.Page)
	w.Header().Set("HX-Push-Url", pushURL)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.ProductsList(listData).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	groupID, ok := parseOptionalGroupID(r.FormValue("group_id"))
	qty, err := parseQuantity(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, "Nieprawidłowa ilość.", http.StatusBadRequest)
		return
	}
	unit := products.Unit(strings.TrimSpace(r.FormValue("unit")))

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	var gid *products.GroupID
	if ok {
		gid = &groupID
	}
	if _, err := s.products.svc.CreateProduct(ctx, products.NewProduct{
		Name:     name,
		GroupID:  gid,
		Quantity: qty,
		Unit:     unit,
	}); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	s.handleProductsPartial(w, r)
}

func (s *Server) handleCreateProductAndRedirect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	groupID, ok := parseOptionalGroupID(r.FormValue("group_id"))
	qty, err := parseQuantity(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, "Nieprawidłowa ilość.", http.StatusBadRequest)
		return
	}
	unit := products.Unit(strings.TrimSpace(r.FormValue("unit")))

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	var gid *products.GroupID
	if ok {
		gid = &groupID
	}
	if _, err := s.products.svc.CreateProduct(ctx, products.NewProduct{
		Name:     name,
		GroupID:  gid,
		Quantity: qty,
		Unit:     unit,
	}); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func (s *Server) handleSetQuantity(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	qty, err := parseQuantity(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, "Nieprawidłowa ilość.", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if err := s.products.svc.SetProductQuantity(ctx, id, qty); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	s.renderProductCard(w, r, id)
}

func (s *Server) handleSetUnit(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	unit := products.Unit(strings.TrimSpace(r.FormValue("unit")))
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if err := s.products.svc.SetProductUnit(ctx, id, unit); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	s.renderProductCard(w, r, id)
}

func (s *Server) handleMarkMissing(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if err := s.products.svc.MarkProductMissing(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	s.renderProductCard(w, r, id)
}

func (s *Server) handleSetGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	groupID, ok := parseOptionalGroupID(r.FormValue("group_id"))
	var gid *products.GroupID
	if ok {
		gid = &groupID
	}

	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if err := s.products.svc.SetProductGroup(ctx, id, gid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	s.renderProductCard(w, r, id)
}

func parseQuantity(v string) (products.Quantity, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(v, 64)
	return products.Quantity(f), err
}

func (s *Server) handleCreateGroup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if _, err := s.products.svc.CreateGroup(ctx, name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	redirect := "/products"
	if r.URL.RawQuery != "" {
		redirect += "?" + r.URL.RawQuery
	}
	w.Header().Set("HX-Redirect", redirect)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleCreateGroupAndRedirect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()
	if _, err := s.products.svc.CreateGroup(ctx, name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.events.Publish(eventProductsList, clientIDFromRequest(r))
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func parsePathProductID(r *http.Request, param string) (products.ProductID, bool) {
	v := r.PathValue(param)
	if v == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(v, 10, 64)
	return products.ProductID(id), err == nil
}

// renderProductCard fetches a single product and renders just its card.
// Used for htmx updates to avoid refreshing the entire product list.
func (s *Server) renderProductCard(w http.ResponseWriter, r *http.Request, id products.ProductID) {
	ctx, cancel := context.WithTimeout(r.Context(), DefaultHandlerTimeout)
	defer cancel()

	// Fetch all products and find the one with matching ID
	productsList, err := s.products.qry.ListProducts(ctx, products.ProductFilter{Limit: 100})
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	var product *products.Product
	for _, p := range productsList {
		if p.ID == id {
			product = &p
			break
		}
	}
	if product == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	// Fetch groups for the card rendering
	groups, err := s.products.qry.ListGroups(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	// Parse current filter state to maintain context
	_, nameQuery, groupNames, formGroupIDs, _ := parseProductsListQuery(r)
	groupIDs := resolveGroupNames(groupNames, groups)
	if len(groupIDs) == 0 {
		groupIDs = formGroupIDs
	}
	onlyMissing := r.URL.Query().Get("missing") == "1"

	listData := views.ProductsListData{
		Groups:           groups,
		OnlyMissing:      onlyMissing,
		NameQuery:        nameQuery,
		SelectedGroupIDs: groupIDs,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.ProductCard(*product, listData).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func parseOptionalGroupID(v string) (products.GroupID, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(v, 10, 64)
	return products.GroupID(id), err == nil
}

func withQuery(r *http.Request, key, value string) *http.Request {
	if value == "" {
		return r
	}
	r2 := r.Clone(r.Context())
	q := r2.URL.Query()
	q.Set(key, value)
	r2.URL.RawQuery = q.Encode()
	return r2
}

func parseProductsListQuery(r *http.Request) (onlyMissing bool, nameQuery string, groupNames []string, groupIDs []products.GroupID, page int64) {
	q := r.URL.Query()
	onlyMissing = q.Get("missing") == "1"
	nameQuery = strings.TrimSpace(q.Get("q"))

	// Parse group names from "groups" parameter (comma-separated) - used in URLs
	if raw := strings.TrimSpace(q.Get("groups")); raw != "" {
		seen := make(map[string]struct{})
		for _, name := range strings.Split(raw, ",") {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if _, dup := seen[name]; dup {
				continue
			}
			seen[name] = struct{}{}
			groupNames = append(groupNames, name)
		}
	}

	// Also parse group_id parameters (used by form checkboxes)
	seenIDs := make(map[products.GroupID]struct{})
	for _, raw := range q["group_id"] {
		id, ok := parseOptionalGroupID(raw)
		if !ok {
			continue
		}
		if _, dup := seenIDs[id]; dup {
			continue
		}
		seenIDs[id] = struct{}{}
		groupIDs = append(groupIDs, id)
	}

	page = 1
	if raw := strings.TrimSpace(q.Get("page")); raw != "" {
		if p, err := strconv.ParseInt(raw, 10, 64); err == nil && p > 0 {
			page = p
		}
	}
	return onlyMissing, nameQuery, groupNames, groupIDs, page
}

// resolveGroupNames converts group names to GroupIDs using the provided groups list.
// Unknown group names are silently ignored.
func resolveGroupNames(groupNames []string, groups []products.Group) []products.GroupID {
	if len(groupNames) == 0 {
		return nil
	}

	// Build name -> ID lookup map
	nameToID := make(map[string]products.GroupID, len(groups))
	for _, g := range groups {
		nameToID[g.Name] = g.ID
	}

	var ids []products.GroupID
	for _, name := range groupNames {
		if id, ok := nameToID[name]; ok {
			ids = append(ids, id)
		}
	}
	return ids
}

// buildProductsPageURL constructs the canonical /products URL for the given filter state.
// Used for HX-Push-Url header to sync browser URL with filter state.
func buildProductsPageURL(onlyMissing bool, nameQuery string, groupIDs []products.GroupID, groups []products.Group, page int64) string {
	values := url.Values{}
	if onlyMissing {
		values.Set("missing", "1")
	}
	nameQuery = strings.TrimSpace(nameQuery)
	if nameQuery != "" {
		values.Set("q", nameQuery)
	}
	if names := products.GroupIDsToNames(groups, groupIDs); len(names) > 0 {
		values.Set("groups", strings.Join(names, ","))
	}
	if page > 1 {
		values.Set("page", strconv.FormatInt(page, 10))
	}
	encoded := values.Encode()
	if encoded == "" {
		return "/products"
	}
	return "/products?" + encoded
}
