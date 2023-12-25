package handler

import (
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

const (
	AuthorizeCookie        = "authorize"
	ContextUser     ctxKey = "context_user"
)

var (
	secureCookie *securecookie.SecureCookie
)

type ctxKey string

type Authorize struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func genId() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

func SetSecureCookie(sc *securecookie.SecureCookie) {
	secureCookie = sc
}
