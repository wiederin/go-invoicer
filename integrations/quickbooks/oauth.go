package quickbooks

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	ProviderID    = "quickbooks"
	InvoiceIDMeta   = "quickbooks_invoice_id"
	qboSandboxAPI = "https://sandbox-quickbooks.api.intuit.com"
	qboProdAPI    = "https://quickbooks.api.intuit.com"
)

// OAuthConfig holds Intuit app credentials.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Sandbox      bool
}

// LoadOAuthConfigFromEnv reads platform QuickBooks OAuth settings.
func LoadOAuthConfigFromEnv(redirectURL string) OAuthConfig {
	sandbox := strings.EqualFold(os.Getenv("GO_INVOICER_QBO_SANDBOX"), "true")
	return OAuthConfig{
		ClientID:     os.Getenv("GO_INVOICER_QBO_CLIENT_ID"),
		ClientSecret: os.Getenv("GO_INVOICER_QBO_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Sandbox:      sandbox,
	}
}

func (c OAuthConfig) configured() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RedirectURL != ""
}

// Configured reports whether server-side QuickBooks OAuth is available.
func (c OAuthConfig) Configured() bool { return c.configured() }

func (c OAuthConfig) OAuthEndpoint() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"com.intuit.quickbooks.accounting"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appcenter.intuit.com/connect/oauth2",
			TokenURL: "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer",
		},
	}
}

// AuthCodeURL builds the Intuit authorization redirect URL.
func (c OAuthConfig) AuthCodeURL(state string) (string, error) {
	if !c.configured() {
		return "", fmt.Errorf("quickbooks: OAuth not configured on server")
	}
	return c.OAuthEndpoint().AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// AccessToken is stored per org after OAuth connect.
type AccessToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	RealmID      string
}

// ExchangeCode exchanges an authorization code for tokens.
func (c OAuthConfig) ExchangeCode(ctx context.Context, code, realmID string) (AccessToken, error) {
	if !c.configured() {
		return AccessToken{}, fmt.Errorf("quickbooks: OAuth not configured on server")
	}
	if realmID == "" {
		return AccessToken{}, fmt.Errorf("quickbooks: realmId missing from callback")
	}
	tok, err := c.OAuthEndpoint().Exchange(ctx, code)
	if err != nil {
		return AccessToken{}, fmt.Errorf("quickbooks: token exchange: %w", err)
	}
	return AccessToken{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.Expiry,
		RealmID:      realmID,
	}, nil
}

// RefreshAccessToken renews an expired access token.
func (c OAuthConfig) RefreshAccessToken(ctx context.Context, refreshToken string) (AccessToken, error) {
	if !c.configured() {
		return AccessToken{}, fmt.Errorf("quickbooks: OAuth not configured on server")
	}
	src := c.OAuthEndpoint().TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	tok, err := src.Token()
	if err != nil {
		return AccessToken{}, fmt.Errorf("quickbooks: refresh token: %w", err)
	}
	return AccessToken{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.Expiry,
	}, nil
}

// APIBase returns the QuickBooks API base URL.
func (c OAuthConfig) APIBase() string {
	if c.Sandbox {
		return qboSandboxAPI
	}
	return qboProdAPI
}
