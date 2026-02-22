package handlers

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/service"
	"github.com/saleh-ghazimoradi/X-Gopher/utils"
	"math"
	"net/http"
	"strconv"
)

type PostHandler struct {
	postService service.PostService
}

func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userId, exists := utils.UserIdFromContext(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid user id", errors.New("invalid user id"))
		return
	}

	var payload dto.CreatePostReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateCreatePostReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid payload")
		return
	}

	post, err := p.postService.CreatePost(r.Context(), userId, &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to create post", err)
		return
	}

	helper.CreatedResponse(w, "Post successfully created", post)
}

func (p *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		helper.BadRequestResponse(w, "Invalid post id", errors.New("invalid post id"))
		return
	}

	post, err := p.postService.GetPostById(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "post not found")
		default:
			helper.InternalServerError(w, "Failed to fetch post", err)
		}
		return
	}

	helper.SuccessResponse(w, "Post fetched successfully", post)
}

func (p *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	if userId == "" {
		helper.BadRequestResponse(w, "user id is required (?id=...)", nil)
		return
	}

	// Pagination from query (with safe defaults)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 10 // default (you can change to 2 if you prefer your Fiber default)
	}

	posts, total, err := p.postService.GetAllPosts(r.Context(), userId, page, limit)
	if err != nil {
		helper.InternalServerError(w, "Failed to fetch feed", err)
		return
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	meta := helper.PaginatedMeta{
		Page:      int64(page),
		Limit:     int64(limit),
		Total:     total,
		TotalPage: totalPages,
	}

	helper.PaginatedSuccessResponse(w, "Posts retrieved successfully", posts, meta)
}

func (p *PostHandler) GetPostsUsersBySearch(w http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	userId, exists := utils.UserIdFromContext(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid user id", errors.New("invalid user id"))
		return
	}

	postId := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if postId == "" {
		helper.BadRequestResponse(w, "Invalid post id", errors.New("invalid post id"))
		return
	}

	var payload dto.UpdatePostReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateUpdatePostReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid payload")
		return
	}

	post, err := p.postService.UpdatePost(r.Context(), postId, userId, &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to update post", err)
		return
	}

	helper.CreatedResponse(w, "Post successfully updated", post)
}

func (p *PostHandler) CommentPost(w http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}
