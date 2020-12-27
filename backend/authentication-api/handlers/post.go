package handlers

import (
	"authentication-api/data"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Signup hanldes user signup requests
func (h *Auth) Signup(w http.ResponseWriter, r *http.Request) {
	h.l.Info("Handling signup request")

	user := r.Context().Value(KeyUser{}).(*data.UserCreate)
	err := h.u.CreateUser(user)

	if err != nil {
		h.l.Error("Something went wrong creating an user in the database ", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}
}

// Signin hanldes user Signin requests
func (h *Auth) Signin(w http.ResponseWriter, r *http.Request) {
	h.l.Info("Handling signin request")

	user := r.Context().Value(KeyUser{}).(*data.UserSignin)
	userdb, err := h.u.GetUserByPhone(user.Phone)

	switch err {
	case nil:
		err = bcrypt.CompareHashAndPassword([]byte(userdb.Password), []byte(user.Password))
		if err != nil {
			h.l.Info("Failed login attempt at", "user", user)
			data.ToJSON(&GenericError{Message: "Can't authenticate user"}, w)
			return
		}
		tokenString, err := h.GenerateToken(user)
		if err != nil {
			h.l.Error("Something went wrong generating token", "error", err)
			data.ToJSON(&GenericError{Message: "Something went wrong generating token"}, w)
			return
		}
		data.ToJSON(&Token{Message: tokenString}, w)
	case data.ErrProductNotFound:
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		h.l.Error("User with phone not found", "id", user.Phone)
		return
	default:
		h.l.Error("Fetching user of id", "id", user.Phone)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
