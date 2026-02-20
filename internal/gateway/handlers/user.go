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

	updatedUser, err := u.userService.UpdateUser(r.Context(), userId, &payload)
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

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
