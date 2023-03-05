package answer

import (
	"context"
	"quizapp/models"
	"quizapp/pkg/types"
)

type Repo interface {
	// Returns created model & nil, if created.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrForbidden, if permission denied.
	// Returns nil & other err else.
	Create(ctx context.Context, answer *models.Answer) (*models.Answer, error)

	// Returns slice & nil, if get smth.
	// Returns empty slice & nil, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetByPoolAnswerId(ctx context.Context, pool_answer_id string, sets types.GetSets) ([]*models.Answer, error)
}
