package oidc

import "net/http"

func RegisterAuthRoutes(mux *http.ServeMux, auth Authenticator) {
	mux.Handle("GET /login", http.HandlerFunc(auth.HandleLogin))
	mux.Handle("GET /oauth2/callback", http.HandlerFunc(auth.HandleCallback))
	mux.Handle("POST /logout", http.HandlerFunc(auth.HandleLogout))
}
