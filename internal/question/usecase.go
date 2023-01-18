package question

import (
	"context"
	"quiz-app/models"
	"quiz-app/pkg/types"
)

type UseCase interface {
	// Returns created model & nil, if created.
	// Returns nil & ErrContentNotFound, if no such form.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrUnauthorized, if user unauthorized.
	// Returns nil & ErrForbidden, if user is not form owner or permission denied.
	// Returns nil & other err else.
	Create(ctx context.Context, model *models.Question) (*models.Question, error)

	// Returns slice & nil, if get smth.
	// Returns empty slice & nil, if get nothing.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & other err else.
	GetByFormId(ctx context.Context, form_id string, sets types.GetSets) ([]*models.Question, error)

	// Returns source model & nil, if updated.
	// Returns nil & ErrContentNotFound, if no such model.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrUnauthorized, if user unauthorized.
	// Returns nil & ErrForbidden, if user is not form owner or permission denied.
	// Returns nil & other err else.
	Update(ctx context.Context, model *models.Question) (*models.Question, error)

	// Returns source model & nil, if deleted.
	// Returns nil & ErrContentNotFound, if no such form or question.
	// Returns nil & ErrInvalidContent, if invalid inputs.
	// Returns nil & ErrUnauthorized, if user unauthorized.
	// Returns nil & ErrForbidden, if user is not form owner or permission denied.
	// Returns nil & other err else.
	Delete(ctx context.Context, id string) error
}
