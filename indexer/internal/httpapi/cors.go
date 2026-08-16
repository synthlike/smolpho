package httpapi

import (
	"fmt"
	"net/http"
	"net/url"
)

// ValidateAllowedOrigins checks that every configured value is a browser
// origin: an HTTP(S) scheme, host, and optional port with no path or query.
func ValidateAllowedOrigins(origins []string) error {
	for _, origin := range origins {
		parsed, err := url.Parse(origin)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
			parsed.Host == "" || parsed.User != nil || parsed.Path != "" ||
			parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("invalid CORS origin %q: expected an origin such as http://localhost:5173", origin)
		}
	}
	return nil
}

// WithCORS allows browser reads from the configured origins. An empty list
// leaves cross-origin browser access disabled.
func WithCORS(next http.Handler, origins []string) http.Handler {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		allowed[origin] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		_, originAllowed := allowed[origin]
		if origin != "" {
			w.Header().Add("Vary", "Origin")
		}
		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		if r.Method != http.MethodOptions || r.Header.Get("Access-Control-Request-Method") == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		if !originAllowed {
			http.Error(w, "CORS origin is not allowed", http.StatusForbidden)
			return
		}
		if r.Header.Get("Access-Control-Request-Method") != http.MethodGet {
			http.Error(w, "CORS method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "600")
		w.WriteHeader(http.StatusNoContent)
	})
}
