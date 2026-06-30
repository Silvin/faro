package auth

import (
	"net/http"
	"time"
)

const sessionCookieName = "faro_session"

func (svc *Service) setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   svc.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (svc *Service) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   svc.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}
