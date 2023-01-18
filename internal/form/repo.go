package form

import (
	"context"
	"quiz-app/models"
	"quiz-app/pkg/types"
)

type Repo interface {
	// Returns created model & nil, if created.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns nil & other err else.
	Create(ctx context.Context, modelBL *models.Form) (*models.Form, error)

	// Returns found model, if get.
	// Returns nil & ErrContentNotFound, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetById(ctx context.Context, id string) (*models.Form, error)

	// Returns slice & nil, if get smth.
	// Returns empty slice & nil, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetByUserId(ctx context.Context, user_id string, sets types.GetSets) ([]*models.Form, error)

	// Returns source model & nil, if updated.
	// Returns nil & ErrContentNotFound, if nothing to update.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns nil & other err else.
	Update(ctx context.Context, modelBL *models.Form) (*models.Form, error)

	// Returns nil, if deleted.
	// Returns ErrContentNotFound, if nothing to delete.
	// Returns ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns other errors else.
	Delete(ctx context.Context, id string) error

	// Returns nil, if user is owner.
	// Returns ErrContentNotFound, if no such form.
	// Returns ErrInvalidContent, if invalid inputs.
	// Returns ErrUnauthorized, if user unauthorized.
	// Returns ErrForbidden, if user is not owner.
	// Returns other err else.
	ValidateIsOwner(ctx context.Context, form_id string) error
}
