package answer

import (
	"context"
	"quiz-app/models"
	"quiz-app/pkg/types"
)

type UseCase interface {
	// Returns created model & nil, if created.
	// Returns empty slice & nil, if get nothing.
	// Returns nil & ErrContentNotFound, if no such form.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrUnauthorized, if user unauthorized.
	// Returns nil & ErrForbidden, if user is not form owner.
	// Returns nil & other err else.
	GetByPoolAnswerId(ctx context.Context, pool_answer_id string, sets types.GetSets) ([]*models.Answer, error)
}
