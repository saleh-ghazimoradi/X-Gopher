package service

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"slices"
)

type UserService interface {
	GetUserById(ctx context.Context, id string) (*dto.UserResp, error)
	GetSuggestedUsers(ctx context.Context, userId string) ([]*dto.UserResp, error)
	UpdateUser(ctx context.Context, id string, input *dto.UpdateUserReq) (*dto.UserResp, error)
	ToggleFollow(ctx context.Context, currentUserId, targetUserId string) (map[string]*dto.UserResp, error)
	DeleteUser(ctx context.Context, id string) error
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

func (u *userService) GetSuggestedUsers(ctx context.Context, userId string) ([]*dto.UserResp, error) {
	mainUser, err := u.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	suggestionSet := make(map[string]struct{})

	for _, followedID := range mainUser.Following {
		followedUser, err := u.userRepository.GetUserById(ctx, followedID)
		if err != nil {
			continue
		}

		for _, id := range followedUser.Following {
			suggestionSet[id] = struct{}{}
		}

		for _, id := range followedUser.Followers {
			suggestionSet[id] = struct{}{}
		}
	}

	delete(suggestionSet, userId)
	for _, id := range mainUser.Following {
		delete(suggestionSet, id)
	}

	ids := make([]string, 0, len(suggestionSet))
	for id := range suggestionSet {
		ids = append(ids, id)
	}

	users, err := u.userRepository.GetUsersByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	resp := make([]*dto.UserResp, 0, len(users))
	for _, user := range users {
		resp = append(resp, u.toUserResp(user))
	}

	return resp, nil
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

func (u *userService) ToggleFollow(ctx context.Context, currentUserId, targetUserId string) (map[string]*dto.UserResp, error) {
	if currentUserId == targetUserId {
		return nil, repository.ErrCannotFollowSelf
	}

	current, err := u.userRepository.GetUserById(ctx, currentUserId)
	if err != nil {
		return nil, err
	}

	target, err := u.userRepository.GetUserById(ctx, targetUserId)
	if err != nil {
		return nil, err
	}

	isFollowing := slices.Contains(target.Followers, currentUserId)

	if isFollowing {
		target.Followers = removeString(target.Followers, currentUserId)
		current.Following = removeString(current.Following, targetUserId)
		if err := u.userRepository.Unfollow(ctx, currentUserId, targetUserId); err != nil {
			return nil, fmt.Errorf("failed to unfollow user: %w", err)
		}
	} else {
		target.Followers = append(target.Followers, currentUserId)
		current.Following = append(current.Following, targetUserId)
		if err := u.userRepository.Follow(ctx, currentUserId, targetUserId); err != nil {
			return nil, fmt.Errorf("failed to follow: %w", err)
		}
	}

	return map[string]*dto.UserResp{
		"target_user":  u.toUserResp(target),
		"current_user": u.toUserResp(current),
	}, nil
}

func (u *userService) DeleteUser(ctx context.Context, id string) error {
	return u.userRepository.DeleteUser(ctx, id)
}

func removeString(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
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
