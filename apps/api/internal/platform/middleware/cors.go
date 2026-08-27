package middleware

import "net/http"

const (
	allowedMethods = "GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS"
	allowedHeaders = "Content-Type"
	exposedHeaders = "Retry-After, X-Request-ID"
)

func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			header := w.Header()
			header.Add("Vary", "Origin")

			if origin != allowedOrigin {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			header.Set("Access-Control-Allow-Origin", allowedOrigin)
			header.Set("Access-Control-Allow-Credentials", "true")
			header.Set("Access-Control-Expose-Headers", exposedHeaders)

			if r.Method == http.MethodOptions {
				header.Add("Vary", "Access-Control-Request-Method")
				header.Add("Vary", "Access-Control-Request-Headers")
				header.Set("Access-Control-Allow-Methods", allowedMethods)
				header.Set("Access-Control-Allow-Headers", allowedHeaders)
				header.Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
