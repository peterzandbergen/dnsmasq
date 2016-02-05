package dnsmasq

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/justinas/alice"
)

type jwtAuthHandler struct {
	h http.Handler
}

func NewJwtAuthHandler(h http.Handler) http.Handler {
	return &jwtAuthHandler{
		h: h,
	}
}

func JwtAuthMiddleWare() alice.Constructor {
	return func(h http.Handler) http.Handler {
		return NewJwtAuthHandler(h)
	}
}

func (h *jwtAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	}
	// Validate token.
	t, err := jwt.ParseFromRequest(r, keyFunc)
	if err != nil || !t.Valid {
		http.Error(w, "Invalid jwt token.", http.StatusUnauthorized)
		return
	}
	// if valid, call the inner handler.
	h.h.ServeHTTP(w, r)
}
