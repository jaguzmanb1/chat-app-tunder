package auth

import (
	"fmt"
	"go-chat/data"
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

// KeyClient usada para el middleware
type KeyClient struct{}

// ValidateToken validates a user request to be signed with a JWT token
func (h *Auth) validateToken(t string, l string) (bool, data.User, error) {
	h.l.Info("[validateToken] Validating token")

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[validateToken] Unexpected signing method: %v", token.Header["alg"])
		}

		return signKey, nil
	})

	if err != nil {
		return false, data.User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		rol := claims["rol"].(float64)
		phone := claims["phone"].(string)

		s := fmt.Sprintf("%.0f", rol)
		if s == l {
			return true, data.User{s, phone}, nil
		}

	}

	return false, data.User{}, err

}
