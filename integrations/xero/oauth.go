package xero

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	ProviderID = "xero"
	InvoiceIDMeta = "xero_invoice_id"
)

// OAuthConfig holds Xero app credentials.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// LoadOAuthConfigFromEnv reads platform Xero OAuth settings.
func LoadOAuthConfigFromEnv(redirectURL string) OAuthConfig {
	return OAuthConfig{
		ClientID:     os.Getenv("GO_INVOICER_XERO_CLIENT_ID"),
		ClientSecret: os.Getenv("GO_INVOICER_XERO_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
	}
}

func (c OAuthConfig) configured() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RedirectURL != ""
}

// Configured reports whether server-side Xero OAuth is available.
func (c OAuthConfig) Configured() bool { return c.configured() }

// OAuthEndpoint returns the Xero OAuth2 config.
func (c OAuthConfig) OAuthEndpoint() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes: []string{
			"openid", "profile", "email",
			"accounting.transactions",
			"accounting.contacts",
			"offline_access",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.xero.com/identity/connect/authorize",
			TokenURL: "https://identity.xero.com/connect/token",
		},
	}
}

// AuthCodeURL builds the authorization redirect URL.
func (c OAuthConfig) AuthCodeURL(state string) (string, error) {
	if !c.configured() {
		return "", fmt.Errorf("xero: OAuth not configured on server")
	}
	return c.OAuthEndpoint().AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// ExchangeCode exchanges an authorization code for tokens and resolves the Xero tenant.
func (c OAuthConfig) ExchangeCode(ctx context.Context, code string) (AccessToken, Tenant, error) {
	if !c.configured() {
		return AccessToken{}, Tenant{}, fmt.Errorf("xero: OAuth not configured on server")
	}
	tok, err := c.OAuthEndpoint().Exchange(ctx, code)
	if err != nil {
		return AccessToken{}, Tenant{}, fmt.Errorf("xero: token exchange: %w", err)
	}
	tenant, err := fetchPrimaryTenant(ctx, tok.AccessToken)
	if err != nil {
		return AccessToken{}, Tenant{}, err
	}
	return AccessToken{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.Expiry,
		TenantID:     tenant.ID,
		TenantName:   tenant.Name,
	}, tenant, nil
}

// AccessToken is stored per org after OAuth connect.
type AccessToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	TenantID     string
	TenantName   string
}

// Tenant is a connected Xero organisation.
type Tenant struct {
	ID   string
	Name string
}

type connection struct {
	TenantID   string `json:"tenantId"`
	TenantName string `json:"tenantName"`
}

func fetchPrimaryTenant(ctx context.Context, accessToken string) (Tenant, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.xero.com/connections", nil)
	if err != nil {
		return Tenant{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Tenant{}, fmt.Errorf("xero: list connections: %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return Tenant{}, fmt.Errorf("xero: list connections: %s", strings.TrimSpace(string(body)))
	}
	var conns []connection
	if err := json.Unmarshal(body, &conns); err != nil {
		return Tenant{}, err
	}
	if len(conns) == 0 {
		return Tenant{}, fmt.Errorf("xero: no Xero organisations connected")
	}
	return Tenant{ID: conns[0].TenantID, Name: conns[0].TenantName}, nil
}

// RefreshAccessToken renews an expired access token.
func (c OAuthConfig) RefreshAccessToken(ctx context.Context, refreshToken string) (AccessToken, error) {
	if !c.configured() {
		return AccessToken{}, fmt.Errorf("xero: OAuth not configured on server")
	}
	src := c.OAuthEndpoint().TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	tok, err := src.Token()
	if err != nil {
		return AccessToken{}, fmt.Errorf("xero: refresh token: %w", err)
	}
	return AccessToken{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.Expiry,
	}, nil
}
