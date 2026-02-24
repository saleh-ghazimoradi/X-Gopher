package service

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"slices"
	"time"
)

type PostService interface {
	CreatePost(ctx context.Context, creatorId string, input *dto.CreatePostReq) (*dto.PostResp, error)
	GetPostById(ctx context.Context, id string) (*dto.PostResp, error)
	GetPostsUsersBySearch(ctx context.Context, query string) (map[string]any, error)
	GetAllPosts(ctx context.Context, userId string, page, limit int) ([]*dto.PostResp, int64, error)
	CommentPost(ctx context.Context, postId, userId string, input *dto.CommentReq) (*dto.PostResp, error)
	LikePost(ctx context.Context, postId, userId string) (*dto.PostResp, error)
	UpdatePost(ctx context.Context, id, userId string, input *dto.UpdatePostReq) (*dto.PostResp, error)
	DeletePost(ctx context.Context, postId, userId string) error
	DeleteComment(ctx context.Context, postId, commentId, userId string) error
}

type postService struct {
	userRepository         repository.UserRepository
	commentRepository      repository.CommentRepository
	postRepository         repository.PostRepository
	notificationRepository repository.NotificationRepository
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

func (p *postService) GetPostsUsersBySearch(ctx context.Context, query string) (map[string]any, error) {
	posts, err := p.postRepository.SearchPosts(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	users, err := p.userRepository.GetUsersBySearch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	postResp := make([]*dto.PostResp, len(posts))
	for i, post := range posts {
		postResp[i] = p.toPostResp(post)
	}

	userResp := make([]*dto.UserResp, len(users))

	for i, u := range users {
		userResp[i] = &dto.UserResp{
			Id:        u.Id,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			ImageUrl:  u.ImageUrl,
			Bio:       u.Bio,
			Followers: u.Followers,
			Following: u.Following,
		}
	}

	return map[string]interface{}{
		"user":  userResp,
		"posts": postResp,
	}, nil
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

func (p *postService) CommentPost(ctx context.Context, postId, userId string, input *dto.CommentReq) (*dto.PostResp, error) {
	comment := &domain.Comment{
		PostId:    postId,
		UserId:    userId,
		Value:     input.Value,
		CreatedAt: time.Now(),
	}

	if err := p.commentRepository.CreateComment(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	if err := p.postRepository.AddComment(ctx, postId, comment.Id); err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	post, err := p.postRepository.GetPostById(ctx, postId)
	if err == nil && post.Creator != userId {
		actor, _ := p.userRepository.GetUserById(ctx, userId)
		notif := &domain.Notification{
			SenderId:   userId,
			ReceiverId: post.Creator,
			TargetId:   postId,
			Details:    actor.FirstName + " " + actor.LastName + " Commented on your Post",
			IsRead:     false,
			CreatedAt:  time.Now(),
			NotificationUser: domain.NotificationUser{
				Name:   actor.FirstName + " " + actor.LastName,
				Avatar: actor.ImageUrl,
			},
		}
		_ = p.notificationRepository.Create(ctx, notif)
	}

	post, _ = p.postRepository.GetPostById(ctx, postId)
	return p.toPostResp(post), nil
}

func (p *postService) LikePost(ctx context.Context, postId, userId string) (*dto.PostResp, error) {
	post, err := p.postRepository.GetPostById(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	isLike := !slices.Contains(post.Likes, userId)

	if isLike {
		post.Likes = append(post.Likes, userId)
		actor, _ := p.userRepository.GetUserById(ctx, userId)
		notif := &domain.Notification{
			SenderId:   userId,
			ReceiverId: post.Creator,
			TargetId:   postId,
			Details:    actor.FirstName + " " + actor.LastName + " Liked your Post",
			IsRead:     false,
			CreatedAt:  time.Now(),
			NotificationUser: domain.NotificationUser{
				Name:   actor.FirstName + " " + actor.LastName,
				Avatar: actor.ImageUrl,
			},
		}
		_ = p.notificationRepository.Create(ctx, notif)
	} else {
		post.Likes = removeString(post.Likes, userId)
	}

	if err := p.postRepository.ToggleLike(ctx, postId, userId); err != nil {
		return nil, fmt.Errorf("failed to toggle like: %w", err)
	}

	return p.toPostResp(post), nil
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

func (p *postService) DeletePost(ctx context.Context, postId, userId string) error {
	post, err := p.postRepository.GetPostById(ctx, postId)
	if err != nil {
		return err
	}

	if post.Creator != userId {
		return fmt.Errorf("post creator does not match")
	}

	return p.postRepository.DeletePost(ctx, postId)
}

func (p *postService) DeleteComment(ctx context.Context, postId, commentId, userId string) error {
	comment, err := p.commentRepository.GetCommentById(ctx, commentId)
	if err != nil {
		return err
	}

	post, err := p.postRepository.GetPostById(ctx, postId)
	if err != nil {
		return err
	}

	if comment.UserId != userId && post.Creator != userId {
		return repository.ErrUnauthorized
	}

	if err := p.commentRepository.DeleteComment(ctx, commentId); err != nil {
		return err
	}

	_ = p.postRepository.RemoveCommentFromPost(ctx, postId, commentId)

	return nil
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

func NewPostService(userRepository repository.UserRepository, commentRepository repository.CommentRepository, postRepository repository.PostRepository, notificationRepository repository.NotificationRepository) PostService {
	return &postService{
		userRepository:         userRepository,
		commentRepository:      commentRepository,
		postRepository:         postRepository,
		notificationRepository: notificationRepository,
	}
}
