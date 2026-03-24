# go-web-infra

Go web infrastructure library containing reusable components for web applications.

This is a library branch. See the main repository for full documentation.

## Packages

### oidc

OIDC (OpenID Connect) authentication middleware and handlers for Go web applications.

#### Features

- OIDC/OAuth2 authentication flow
- Session management with secure cookies
- Middleware for protecting routes
- Support for disabled auth mode (development)
- Admin claim extraction from OIDC tokens

#### Usage

```go
import "github.com/jborkows/k3sfun/go-web-infra/oidc"

// Create authenticator
cfg := oidc.AuthenticationConfig{
    OIDCIssuer:       "https://auth.example.com",
    OIDCClientID:     "my-client-id",
    OIDCClientSecret: "my-client-secret",
    OIDCRedirectURL:  "https://myapp.example.com/oauth2/callback",
    AuthDisabled:     false, // Set to true for development
}

auth, err := oidc.New(cfg)
if err != nil {
    log.Fatal(err)
}

// Use middleware
http.Handle("/", auth.Middleware(handler))

// Handle auth routes
http.HandleFunc("/login", auth.HandleLogin)
http.HandleFunc("/oauth2/callback", auth.HandleCallback)
http.HandleFunc("/logout", auth.HandleLogout)

// Get current user
user, ok := auth.CurrentUser(r)
if ok {
    fmt.Println(user.Email, user.Name, user.Admin)
}
```

#### Configuration

| Field | Description | Required |
|-------|-------------|----------|
| `OIDCIssuer` | OIDC provider URL (e.g., `https://accounts.google.com`) | Yes (unless AuthDisabled) |
| `OIDCClientID` | Client ID from OIDC provider | Yes (unless AuthDisabled) |
| `OIDCClientSecret` | Client secret from OIDC provider | Yes (unless AuthDisabled) |
| `OIDCRedirectURL` | Callback URL (e.g., `https://app.example.com/oauth2/callback`) | Yes (unless AuthDisabled) |
| `AuthDisabled` | Disable authentication (returns dev user) | No (default: false) |

#### User struct

```go
type User struct {
    Subject string  // OIDC subject identifier
    Email   string  // User email
    Name    string  // User display name
    Admin   bool    // Admin flag from claims
}
```

## License

MIT
