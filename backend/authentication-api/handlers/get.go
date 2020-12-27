package handlers

import (
	"authentication-api/data"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken a token
func (h *Auth) GenerateToken(user *data.UserSignin) (string, error) {
	h.l.Info("Generating token for user", "user", user)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["phone"] = user.Phone
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["rol"] = user.Rol

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
