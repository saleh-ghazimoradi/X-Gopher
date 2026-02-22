package routes

import (
	"github.com/julienschmidt/httprouter"
	_ "github.com/saleh-ghazimoradi/X-Gopher/docs"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Register struct {
	authRoute   *AuthRoute
	userRoute   *UserRoute
	postRoute   *PostRoute
	middlewares *middlewares.Middleware
}

type Options func(*Register)

func WithAuthRoute(authRoute *AuthRoute) Options {
	return func(r *Register) {
		r.authRoute = authRoute
	}
}

func WithUserRoute(userRoute *UserRoute) Options {
	return func(r *Register) {
		r.userRoute = userRoute
	}
}

func WithPostRoute(postRoute *PostRoute) Options {
	return func(r *Register) {
		r.postRoute = postRoute
	}
}

func WithMiddlewares(middlewares *middlewares.Middleware) Options {
	return func(r *Register) {
		r.middlewares = middlewares
	}
}

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

	r.authRoute.AuthRoutes(router)
	r.userRoute.UserRoutes(router)
	r.postRoute.PostRoutes(router)
	return r.middlewares.Recover(r.middlewares.Logging(r.middlewares.CORS(r.middlewares.RateLimit(router))))
}

func NewRegister(opts ...Options) *Register {
	register := &Register{}
	for _, opt := range opts {
		opt(register)
	}
	return register
}
