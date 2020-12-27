package handlers

import (
	"authentication-api/data"
	"os"

	"github.com/hashicorp/go-hclog"
)

var signKey = []byte(os.Getenv("secretAuthKey"))

// Auth describes a Auth http handler object
type Auth struct {
	l hclog.Logger
	u *data.UserService
	v *data.Validation
}

// KeyUser used for the middleware to pass data trought request context
type KeyUser struct{}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// Token is the token structure generated by the ser server
type Token struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// New creates a new instance of a auth handler
func New(l hclog.Logger, u *data.UserService, v *data.Validation) *Auth {
	return &Auth{l, u, v}
}