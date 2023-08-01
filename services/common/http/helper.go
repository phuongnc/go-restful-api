package http

import (
	"net/http"
)

func HandlerFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		applySecurityProtectionForResponseHeader(w)
		handler.ServeHTTP(w, r)
	}
}

func applySecurityProtectionForResponseHeader(w http.ResponseWriter) {
	w.Header().Set("Strict-Transport-Security", "max-age=7776000; includeSubDomains")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	//TODO: In order to avoid any XSS/CSRF attack in server side rendering / static files serve in the future, we need to strictly set the valid sources for all static/ dynamic resources.
	// For now, I set it to * to accept resources can be loaded from any external sources
	w.Header().Set("Content-Security-Policy", "default-src 'self' *;connect-src 'self' *;img-src 'self' *;font-src 'self' *;script-src 'self' * 'unsafe-inline';style-src 'self' * 'unsafe-inline'")
}
