package web

import (
	"testing"

	"github.com/jborkows/k3sfun/go-web-infra/oidc"
	"shopping/internal/infrastructure/config"
)

func TestServerIsAdmin_AdminClaim(t *testing.T) {
	s := &Server{cfg: config.Config{}}

	if !s.isAdmin(&oidc.User{Admin: true}) {
		t.Fatalf("expected user to be admin via admin claim")
	}
}

func TestServerIsAdmin_EmailFallback(t *testing.T) {
	s := &Server{cfg: config.Config{AdminEmails: []string{"admin@example.com"}}}

	if !s.isAdmin(&oidc.User{Email: "admin@example.com"}) {
		t.Fatalf("expected user to be admin via email match")
	}
}

func TestServerIsAdmin_NoConfig(t *testing.T) {
	s := &Server{cfg: config.Config{}}

	if s.isAdmin(&oidc.User{Email: "admin@example.com"}) {
		t.Fatalf("expected user to not be admin when no admin config is set")
	}
}
