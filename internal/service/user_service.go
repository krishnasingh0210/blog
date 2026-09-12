package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/krishnasingh0210/blog/internal/domain"
	"github.com/krishnasingh0210/blog/pkg/jwt"
)

type userService struct {
	repo       domain.UserRepository
	jwtManager *jwt.Manager
}

func NewUserService(repo domain.UserRepository, jwtManager *jwt.Manager) domain.UserService {
	return &userService{repo: repo, jwtManager: jwtManager}
}

func (s *userService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         domain.RoleAdmin,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err // repo already translates duplicate email -> domain.ErrDuplicateData
	}
	return u, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// deliberately return the SAME error for "no such user" and "wrong password" below —
			// never let a login endpoint reveal whether an email exists in your system
			return "", domain.ErrInvalidCreds
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCreds
	}

	return s.jwtManager.Generate(u.ID.String(), string(u.Role))
}