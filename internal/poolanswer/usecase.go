package poolanswer

import (
	"context"
	"quizapp/models"
	"quizapp/pkg/types"
)

type UseCase interface {
	// Returns created model & nil, if created.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns nil & other err else.
	Create(ctx context.Context, pool_answer *models.PoolAnswer, answers []*models.Answer) (*models.PoolAnswer, []*models.Answer, error)

	// Returns created model & nil, if created.
	// Returns empty slice & nil, if get nothing.
	// Returns nil & ErrContentNotFound, if no such form.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrUnauthorized, if user unauthorized.
	// Returns nil & ErrForbidden, if user is not form owner.
	// Returns nil & other err else.
	GetByFormId(ctx context.Context, form_id string, sets types.GetSets) ([]*models.PoolAnswer, error)

	// Returns found model, if get.
	// Returns nil & ErrContentNotFound, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetById(ctx context.Context, id string) (*models.PoolAnswer, error)
}
