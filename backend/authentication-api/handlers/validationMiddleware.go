package handlers

import (
	"authentication-api/data"
	"context"
	"net/http"
)

//MiddlewareValidateUser  verificacion para los request
func (h *Auth) MiddlewareValidateUser(next http.Handler) http.Handler {
	h.l.Info("[MiddlewareValidateUser] Handling validator middleware request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		user := &data.UserCreate{}

		err := data.FromJSON(user, r.Body)
		if err != nil {
			h.l.Error("[MiddlewareValidateUser] Deserializing product", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
		h.l.Debug("[MiddlewareValidateUser] Serialized user", "user", user)
		errs := h.v.Validate(user)
		if len(errs) != 0 {
			h.l.Error("[MiddlewareValidateUser] Validating user", "errors:", errs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

//MiddlewareValidateUserSignin verificacion para los request
func (h *Auth) MiddlewareValidateUserSignin(next http.Handler) http.Handler {
	h.l.Info("[MiddlewareValidateUserSignin] Handling validator middleware request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		user := &data.UserSignin{}

		err := data.FromJSON(user, r.Body)
		if err != nil {
			h.l.Error("[MiddlewareValidateUserSignin] Deserializing product", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
		h.l.Debug("[MiddlewareValidateUserSignin] Serialized user", "phone", user.Phone)
		errs := h.v.Validate(user)
		if len(errs) != 0 {
			h.l.Error("[MiddlewareValidateUserSignin] Validating user", "errors:", errs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
