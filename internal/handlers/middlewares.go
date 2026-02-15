package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyHash, err := getAPIKeyHash(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var (
			id     string
			domainId string
			domain string
		)
		err = s.pool.QueryRow(`
			SELECT k.id, d.id AS domain_id, d.domain
			FROM api_keys AS k
			INNER JOIN domains AS d ON (k.domain_id=d.id)
			WHERE key_hash=$1
				AND k.is_active=true
				AND d.is_active=true
		`,
			apiKeyHash,
		).Scan(&id, &domainId, &domain)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "api_key", id)
		ctx = context.WithValue(ctx, "domain_id", domainId)
		ctx = context.WithValue(ctx, "domain", domain)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAPIKeyHash(authorization string) (string, error) {
	bearerParts := strings.Split(authorization, " ")
	if len(bearerParts) != 2 {
		return "", fmt.Errorf("No API key provided")
	}
	apiKeyParts := bearerParts[1]
	if !strings.HasPrefix(apiKeyParts, "sk_live_") {
		return "", fmt.Errorf("API key is malformed")
	}
	apiKeyHash := sha256.Sum256([]byte(strings.TrimPrefix(apiKeyParts, "sk_live_")))

	return fmt.Sprintf("%x", apiKeyHash), nil
}
