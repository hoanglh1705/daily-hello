package services

import (
	"context"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	// Check duplicate email
	existing, _ := s.repo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, appErrors.ErrEmailExists
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	user := &models.User{
		Name:     req.Name,
		Code:     req.Code,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashed),
		Role:     req.Role,
		BranchID: req.BranchID,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, appErrors.ErrInternal
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrNotFound
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id uint, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.BranchID != nil {
		user.BranchID = req.BranchID
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, appErrors.ErrInternal
	}

	return user, nil
}

func (s *UserService) List(ctx context.Context, q models.UserListQuery) (*models.PaginatedResponse, error) {
	items, total, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return &models.PaginatedResponse{
		Items: items,
		Meta: models.PaginationMeta{
			Page:  q.GetPage(),
			Limit: q.GetLimit(),
			Total: total,
		},
	}, nil
}
