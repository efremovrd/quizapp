package form

import (
	"context"
	"quiz-app/models"
	"quiz-app/pkg/types"
)

type UseCase interface {
	// Returns created model & nil, if created.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns nil & other err else.
	Create(ctx context.Context, model *models.Form) (*models.Form, error)

	// Returns slice & nil, if get smth.
	// Returns empty slice & nil, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetByUserId(ctx context.Context, user_id string, sets types.GetSets) ([]*models.Form, error)

	// Returns found models & nil, if get.
	// Returns nil & ErrContentNotFound, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetById(ctx context.Context, id string) (*models.Form, error)

	// Returns source model & nil, if updated.
	// Returns nil & ErrContentNotFound, if nothing to update.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrUnauthorized, if user not set in context.
	// Returns nil & ErrForbidden, if user is not an owner or permission denied.
	// Returns nil & other err else.
	Update(ctx context.Context, model *models.Form) (*models.Form, error)

	// Returns nil, if deleted.
	// Returns ErrContentNotFound, if nothing to delete.
	// Returns ErrInvalidContent, if invalid inputs.
	// Returns ErrUnauthorized, if user not set in context.
	// Returns ErrForbidden, if user is not an owner or permission denied.
	// Returns other errors else.
	Delete(ctx context.Context, id string) error
}
