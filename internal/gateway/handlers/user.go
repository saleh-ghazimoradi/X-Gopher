package handlers

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/service"
	"github.com/saleh-ghazimoradi/X-Gopher/utils"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func (u *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		helper.BadRequestResponse(w, "Invalid given user id", errors.New("invalid user id"))
		return
	}

	user, err := u.userService.GetUserById(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "User not found")
		default:
			helper.InternalServerError(w, "Internal server error", err)
		}
		return
	}

	helper.SuccessResponse(w, "user successfully retrieved", user)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, exists := utils.UserIdFromContext(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid given user id", errors.New("invalid user id"))
		return
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" || id != userId {
		helper.UnauthorizedResponse(w, "You can only update your own profile")
		return
	}

	var payload dto.UpdateUserReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateUpdateUserReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid given payload")
		return
	}

	updatedUser, err := u.userService.UpdateUser(r.Context(), id, &payload)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "User not found")
		default:
			helper.InternalServerError(w, "Internal server error", err)
		}
		return
	}

	helper.SuccessResponse(w, "user successfully updated", updatedUser)
}

func (u *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	targetId := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if targetId == "" {
		helper.BadRequestResponse(w, "Invalid given user id", errors.New("invalid user id"))
		return
	}

	currentUserId, exists := utils.UserIdFromContext(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid user id from token", errors.New("user id not found in context"))
		return
	}

	resp, err := u.userService.ToggleFollow(r.Context(), currentUserId, targetId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "User not found")
		case errors.Is(err, repository.ErrCannotFollowSelf):
			helper.BadRequestResponse(w, "Cannot follow yourself", err)
		default:
			helper.InternalServerError(w, "Failed to toggle follow", err)
		}
		return
	}

	helper.SuccessResponse(w, "Follow status toggled successfully", resp)
}

func (u *UserHandler) GetSuggestedUsers(w http.ResponseWriter, r *http.Request) {
	userId, exists := utils.UserIdFromContext(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid user id", errors.New("user id not found in context"))
		return
	}

	users, err := u.userService.GetSuggestedUsers(r.Context(), userId)
	if err != nil {
		helper.InternalServerError(w, "Failed to get suggested users", err)
		return
	}

	helper.SuccessResponse(w, "Suggested users retrieved successfully", users)
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, exists := utils.UserIdFromContext(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid given user id", errors.New("invalid user id"))
		return
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" || id != userId {
		helper.UnauthorizedResponse(w, "You can only update your own profile")
		return
	}

	if err := u.userService.DeleteUser(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "User not found")
		default:
			helper.InternalServerError(w, "Failed to delete user", err)
		}
		return
	}

	helper.SuccessResponse(w, "User deleted successfully", nil)
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
