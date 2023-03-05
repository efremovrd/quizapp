package auth

import (
	"context"
	"quizapp/models"
)

type UseCase interface {
	// Returns signed up model & nil, if signed up.
	// Returns nil & ErrLoginExists, if login exists.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns nil & other err else.
	SignUp(ctx context.Context, user *models.User) (*models.User, error)

	// Returns token & nil, if successful signing in.
	// Returns nil & ErrUnathorized, if no such user.
	// Returns nil & ErrInvalidPassword, if invalid password.
	// Returns nil & other err else.
	SignIn(ctx context.Context, user *models.User) (*string, error)

	// Returns found model & nil, if get.
	// Returns nil & ErrContentNotFound, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetById(ctx context.Context, id string) (*models.User, error)

	// Returns user model & nil, if parsed.
	// Returns nil & ErrInvalidAccessToken else.
	ParseToken(ctx context.Context, token string) (*models.User, error)
}
