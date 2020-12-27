package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-hclog"
)

var signKey = []byte(os.Getenv("secretAuthKey"))

// Auth define a object of auth to validate requests tokens
type Auth struct {
	l hclog.Logger
}

// New creates a new auth validator instance
func New(l hclog.Logger) *Auth {
	return &Auth{l}
}

// ValidateToken validates a user request to be signed with a JWT token
func (h *Auth) validateToken(t string) (bool, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			h.l.Error("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["authorized"], claims["user"], claims["exp"])
		return true, nil
	}

	return false, err

}

// MiddlewareTokenValidation validates resquests tokens
func (h *Auth) MiddlewareTokenValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			tv, err := h.validateToken(r.Header["Authorization"][0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.l.Error("Error parsing or validating request token", "error", err)

				return
			}

			if tv {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		h.l.Info("User request not autorized")
		fmt.Fprintf(w, "User request not autorized")
	})
}
