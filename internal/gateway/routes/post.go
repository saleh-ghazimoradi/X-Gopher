package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/middlewares"
	"net/http"
)

type PostRoute struct {
	middlewares *middlewares.Middleware
	postHandler *handlers.PostHandler
}

func (p *PostRoute) PostRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/v1/post/:id", p.postHandler.GetPost)
	router.HandlerFunc(http.MethodGet, "/v1/post", p.postHandler.GetAllPosts)
	router.HandlerFunc(http.MethodGet, "/v1/postSearch", p.postHandler.GetPostsUsersBySearch)

	router.Handler(http.MethodPost, "/v1/post", p.wrapAuth(p.postHandler.CreatePost))
	router.Handler(http.MethodPatch, "/v1/post/:id", p.wrapAuth(p.postHandler.UpdatePost))
	router.Handler(http.MethodPost, "/v1/post/:id/comment", p.wrapAuth(p.postHandler.CommentPost))
	router.Handler(http.MethodPatch, "/v1/post/:id/like", p.wrapAuth(p.postHandler.LikePost))
	router.Handler(http.MethodDelete, "/v1/post/:id", p.wrapAuth(p.postHandler.DeletePost))
	router.Handler(http.MethodDelete, "/v1/comments/:postId/comments/:commentId", p.wrapAuth(p.postHandler.DeleteComment))
}

func (p *PostRoute) wrapAuth(handler http.HandlerFunc) http.Handler {
	return p.middlewares.Authenticate(handler)
}

func NewPostRoute(middlewares *middlewares.Middleware, postHandler *handlers.PostHandler) *PostRoute {
	return &PostRoute{
		middlewares: middlewares,
		postHandler: postHandler,
	}
}
