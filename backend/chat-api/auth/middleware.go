package auth

import (
	"fmt"
	"go-chat/data"
	"net/http"
	"time"

	"github.com/gorilla/context"

	"github.com/gorilla/websocket"
)

// MiddlewareTokenValidationRol0 validates resquests tokens
func (h *Auth) MiddlewareTokenValidationRol0(next http.Handler) http.Handler {
	h.l.Info("[MiddlewareTokenValidationRol0] Handling validator middleware request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			tv, _, err := h.validateToken(r.Header["Authorization"][0], "0")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.l.Error("[MiddlewareTokenValidationRol1] Error parsing or validating request token", "error", err, "endpoint", r.URL)
				return
			}

			if tv {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		h.l.Info("[MiddlewareTokenValidationRol0] User request not autorized")
		fmt.Fprintf(w, "User request not autorized")
	})
}

// MiddlewareTokenValidationRol1 validates resquests tokens
func (h *Auth) MiddlewareTokenValidationRol1(next http.Handler) http.Handler {
	h.l.Info("[MiddlewareTokenValidationRol1] Handling validator middleware request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			tv, _, err := h.validateToken(r.Header["Authorization"][0], "1")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.l.Error("[MiddlewareTokenValidationRol1] Error parsing or validating request token", "error", err, "endpoint", r.URL)

				return
			}

			if tv {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		h.l.Info("[MiddlewareTokenValidationRol1] User request not autorized")
		fmt.Fprintf(w, "User request not autorized")
	})
}

var upgrader = websocket.Upgrader{}

// MiddlewareTokenValidationSocket validates resquests tokens
func (h *Auth) MiddlewareTokenValidationSocket(next http.Handler) http.Handler {
	h.l.Info("[MiddlewareTokenValidationSocket] Handling validator middleware request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg data.Message
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		ws.SetReadDeadline(time.Now().Add(5 * time.Second))
		err = ws.ReadJSON(&msg)
		defer ws.Close()

		switch err {
		case nil:
			tv, user, _ := h.validateToken(msg.Message, "0")
			if tv {
				context.Set(r, "ws", ws)
				context.Set(r, "us", user)

				next.ServeHTTP(w, r)
				return
			}
			h.l.Info("[MiddlewareTokenValidationSocket] User request not autorized")
			return
		default:
			h.l.Error("MiddlewareTokenValidationSocket] Could not make first websocket handshake", "error", err)
			return
		}

	})
}
