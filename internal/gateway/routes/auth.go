package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/handlers"
	"net/http"
)

type AuthRoute struct {
	authHandler *handlers.AuthHandler
}

func (a *AuthRoute) AuthRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/auth/signup", a.authHandler.Register)
	router.HandlerFunc(http.MethodPost, "/v1/auth/login", a.authHandler.Login)
	router.HandlerFunc(http.MethodPost, "/v1/auth/refresh", a.authHandler.RefreshToken)
	router.HandlerFunc(http.MethodPost, "/v1/auth/logout", a.authHandler.Logout)
}

func NewAuthRoute(authHandler *handlers.AuthHandler) *AuthRoute {
	return &AuthRoute{
		authHandler: authHandler,
	}
}
