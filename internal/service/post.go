package service

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"time"
)

type PostService interface {
	CreatePost(ctx context.Context, creatorId string, input *dto.CreatePostReq) (*dto.PostResp, error)
	GetPostById(ctx context.Context, id string) (*dto.PostResp, error)
	GetAllPosts(ctx context.Context, userId string, page, limit int) ([]*dto.PostResp, int64, error)
	UpdatePost(ctx context.Context, id, userId string, input *dto.UpdatePostReq) (*dto.PostResp, error)
}

type postService struct {
	userRepository repository.UserRepository
	postRepository repository.PostRepository
}

func (p *postService) CreatePost(ctx context.Context, creatorId string, input *dto.CreatePostReq) (*dto.PostResp, error) {
	user, err := p.userRepository.GetUserById(ctx, creatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get creator: %w", err)
	}

	post := &domain.Post{
		Creator:      creatorId,
		Title:        input.Title,
		Message:      input.Message,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		SelectedFile: input.SelectedFile,
		Likes:        make([]string, 0),
		Comments:     make([]string, 0),
		CreatedAt:    time.Now(),
	}

	if err := p.postRepository.CreatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return p.toPostResp(post), nil
}

func (p *postService) GetPostById(ctx context.Context, id string) (*dto.PostResp, error) {
	post, err := p.postRepository.GetPostById(ctx, id)
	if err != nil {
		return nil, err
	}

	return p.toPostResp(post), nil
}

func (p *postService) GetAllPosts(ctx context.Context, userId string, page, limit int) ([]*dto.PostResp, int64, error) {
	mainUser, err := p.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	// Build feed: user + everyone they follow
	feedIds := append([]string{userId}, mainUser.Following...)

	posts, total, err := p.postRepository.GetFeedPosts(ctx, feedIds, page, limit)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]*dto.PostResp, len(posts))
	for i, post := range posts {
		resp[i] = p.toPostResp(post)
	}

	return resp, total, nil
}

func (p *postService) UpdatePost(ctx context.Context, id, userId string, input *dto.UpdatePostReq) (*dto.PostResp, error) {
	post, err := p.postRepository.GetPostById(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.Creator != userId {
		return nil, fmt.Errorf("post creator does not match")
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Message != nil {
		post.Message = *input.Message
	}

	if input.SelectedFile != nil {
		post.SelectedFile = *input.SelectedFile
	}

	if err := p.postRepository.UpdatePost(ctx, post); err != nil {
		return nil, err
	}

	return p.toPostResp(post), nil
}

func (p *postService) toPostResp(input *domain.Post) *dto.PostResp {
	return &dto.PostResp{
		Id:           input.Id,
		Creator:      input.Creator,
		Title:        input.Title,
		Message:      input.Message,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		SelectedFile: input.SelectedFile,
		Likes:        input.Likes,
		Comments:     input.Comments,
		CreatedAt:    input.CreatedAt,
	}
}

func NewPostService(userRepository repository.UserRepository, postRepository repository.PostRepository) PostService {
	return &postService{
		userRepository: userRepository,
		postRepository: postRepository,
	}
}
