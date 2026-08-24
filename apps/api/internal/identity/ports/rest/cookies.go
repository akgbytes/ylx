package rest

import (
	"net/http"
	"time"

	"github.com/akgbytes/ylx/internal/identity/app"
	"github.com/akgbytes/ylx/internal/platform/config"
)

type Cookies struct {
	accessName    string
	refreshName   string
	accessMaxAge  time.Duration
	refreshMaxAge time.Duration
	secure        bool
}

func NewCookies(cfg *config.AuthConfig, secure bool) *Cookies {
	return &Cookies{
		accessName:    cfg.AccessTokenName,
		refreshName:   cfg.RefreshTokenName,
		accessMaxAge:  cfg.AccessTokenExpiry,
		refreshMaxAge: cfg.RefreshTokenExpiry,
		secure:        secure,
	}
}

func (c *Cookies) AccessToken(r *http.Request) string {
	return cookieValue(r, c.accessName)
}

func (c *Cookies) RefreshToken(r *http.Request) string {
	return cookieValue(r, c.refreshName)
}

func (c *Cookies) Set(w http.ResponseWriter, tokens app.Tokens) {
	c.write(w, c.accessName, tokens.AccessToken, tokens.AccessTokenExpiresAt, int(c.accessMaxAge.Seconds()))
	c.write(w, c.refreshName, tokens.RefreshToken, tokens.RefreshTokenExpiresAt, int(c.refreshMaxAge.Seconds()))
}

func (c *Cookies) Clear(w http.ResponseWriter) {
	expired := time.Unix(1, 0)

	c.write(w, c.accessName, "", expired, -1)
	c.write(w, c.refreshName, "", expired, -1)
}

func (c *Cookies) write(w http.ResponseWriter, name, value string, expires time.Time, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func cookieValue(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
