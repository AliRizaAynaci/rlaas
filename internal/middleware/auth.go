package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string
		if c, err := r.Cookie("session_token"); err == nil {
			tokenStr = c.Value
		}
		if tokenStr == "" {
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenStr = strings.TrimPrefix(auth, "Bearer ")
			}
		}
		if tokenStr == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !tok.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims := tok.Claims.(jwt.MapClaims)
		uid := int(claims["user_id"].(float64))

		ctx := context.WithValue(r.Context(), UserIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FromContext(ctx context.Context) int {
	if v, ok := ctx.Value(UserIDKey).(int); ok {
		return v
	}
	return 0
}
