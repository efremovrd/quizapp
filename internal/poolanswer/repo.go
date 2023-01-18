package poolanswer

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
	Create(ctx context.Context, pool_answer *models.PoolAnswer) (*models.PoolAnswer, error)

	// Returns slice & nil, if get smth.
	// Returns empty slice & nil, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetByFormId(ctx context.Context, form_id string, sets types.GetSets) ([]*models.PoolAnswer, error)

	// Returns nil, if deleted.
	// Returns ErrContentNotFound, if nothing to delete.
	// Returns ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns other errors else.
	Delete(ctx context.Context, id string) error

	// Returns found model, if get.
	// Returns nil & ErrContentNotFound, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetById(ctx context.Context, id string) (*models.PoolAnswer, error)
}