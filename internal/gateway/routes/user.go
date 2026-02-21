package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/middlewares"
	"net/http"
)

type UserRoute struct {
	middlewares *middlewares.Middleware
	userHandler *handlers.UserHandler
}

func (u *UserRoute) UserRoutes(router *httprouter.Router) {
	// Public Routes
	router.HandlerFunc(http.MethodGet, "/v1/user/:id", u.userHandler.GetUserById)

	// Protected Routes
	router.Handler(http.MethodPatch, "/v1/user/:id", u.wrapAuth(u.userHandler.UpdateUser))
	router.Handler(http.MethodPatch, "/v1/user/:id/following", u.wrapAuth(u.userHandler.FollowUser))
}

func (u *UserRoute) wrapAuth(handler http.HandlerFunc) http.Handler {
	return u.middlewares.Authenticate(handler)
}

func NewUserRoute(middlewares *middlewares.Middleware, userHandler *handlers.UserHandler) *UserRoute {
	return &UserRoute{
		middlewares: middlewares,
		userHandler: userHandler,
	}
}
