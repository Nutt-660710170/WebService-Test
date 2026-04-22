package handler

import (
	"net/http"

	"gorm.io/gorm"
	"oil/internal/domain"
)

// APIKeyAuthMiddleware checks for a valid API_KEY in the Authorization header
func APIKeyAuthMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		apiKey := header[len(prefix):]

		var key domain.APIKey
		err := db.Where("key = ? AND status = ?", apiKey, "active").First(&key).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if key.ExpiredAt != nil && key.ExpiredAt.Before(r.Context().Value("now").(func() time.Time)()) {
			http.Error(w, "API key expired", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
