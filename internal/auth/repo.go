package auth

import (
	"context"
	"quiz-app/models"
)

type Repo interface {
	// Returns created model & nil, if created.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns nil & other err else.
	Create(ctx context.Context, user *models.User) (*models.User, error)

	// Returns found model & nil, if get.
	// Returns nil & ErrContentNotFound, if get nothing.
	// Returns nil & other err else.
	GetByLogin(ctx context.Context, login string) (*models.User, error)

	// Returns found model & nil, if get.
	// Returns nil & ErrContentNotFound, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetById(ctx context.Context, id string) (*models.User, error)
}
