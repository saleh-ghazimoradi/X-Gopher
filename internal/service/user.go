package service

import (
	"context"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
)

type UserService interface {
	GetUserById(ctx context.Context, id string) (*dto.UserResp, error)
	UpdateUser(ctx context.Context, id string, input *dto.UpdateUserReq) (*dto.UserResp, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) GetUserById(ctx context.Context, id string) (*dto.UserResp, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.toUserResp(user), nil
}

func (u *userService) UpdateUser(ctx context.Context, id string, input *dto.UpdateUserReq) (*dto.UserResp, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	if input.ImageUrl != nil {
		user.ImageUrl = *input.ImageUrl
	}

	if input.Bio != nil {
		user.Bio = *input.Bio
	}

	if err := u.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return u.toUserResp(user), nil
}

func (u *userService) toUserResp(input *domain.User) *dto.UserResp {
	return &dto.UserResp{
		Id:        input.Id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		ImageUrl:  input.ImageUrl,
		Bio:       input.Bio,
		Followers: input.Followers,
		Following: input.Following,
	}
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
