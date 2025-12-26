package web

import (
	"net/http"
	"strings"
)

func clientIDFromRequest(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("X-Client-ID"))
}
