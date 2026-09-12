package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/krishnasingh0210/blog/internal/domain"
	"github.com/krishnasingh0210/blog/internal/repository/sqlcgen"
)

type userRepository struct {
	queries *sqlcgen.Queries
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{queries: sqlcgen.New(db)}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	row, err := r.queries.CreateUser(ctx, sqlcgen.CreateUserParams{
		Name:     u.Name,
		Email:    u.Email,
		PasswordHash: u.PasswordHash,
		Role:      string(u.Role),
	})
	if err != nil {
		if isUniqueViolation(err){
			return domain.ErrDuplicateData
		}
		return err
	}
	*u = toDomainUser(row)
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u := toDomainUser(row)
	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u := toDomainUser(row)
	return &u, nil
}


func toDomainUser (row sqlcgen.User) domain.User {
	return domain.User{
		ID:     row.ID,
		Name:   row.Name,
		Email:  row.Email,
		PasswordHash: row.PasswordHash,
		Role:    domain.Role(row.Role),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}