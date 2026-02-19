package routes

import (
	"github.com/julienschmidt/httprouter"
	_ "github.com/saleh-ghazimoradi/X-Gopher/docs"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Register struct{}

type Options func(*Register)

func (r *Register) RegisterRoutes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(helper.HTTPRouterNotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(helper.HTTPRouterMethodNotAllowedResponse)

	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))

	router.Handler("GET", "/docs/*filepath",
		http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))),
	)

	router.Handler("GET", "/api-docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/rapidoc.html")
	}))

	return router
}

func NewRegister(opts ...Options) *Register {
	register := &Register{}
	for _, opt := range opts {
		opt(register)
	}
	return register
}
