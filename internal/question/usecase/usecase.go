package usecase

import (
	"context"
	"quizapp/internal/form"
	"quizapp/internal/question"
	"quizapp/models"
	"quizapp/pkg/types"
)

type questionUseCase struct {
	qRepo question.Repo
	fRepo form.Repo
}

func NewQuestionUseCase(qRepo question.Repo, fRepo form.Repo) question.UseCase {
	return &questionUseCase{
		qRepo: qRepo,
		fRepo: fRepo,
	}
}

func (q *questionUseCase) Create(ctx context.Context, model *models.Question) (*models.Question, error) {
	err := q.fRepo.ValidateIsOwner(ctx, model.Form_id)
	if err != nil {
		return nil, err
	}

	return q.qRepo.Create(ctx, model)
}

func (q *questionUseCase) GetByFormId(ctx context.Context, form_id string, sets types.GetSets) ([]*models.Question, error) {
	return q.qRepo.GetByFormId(ctx, form_id, sets)
}

func (q *questionUseCase) Update(ctx context.Context, model *models.Question) (*models.Question, error) {
	err := q.fRepo.ValidateIsOwner(ctx, model.Form_id)
	if err != nil {
		return nil, err
	}

	_, err = q.qRepo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (q *questionUseCase) Delete(ctx context.Context, id string) error {
	foundquestion, err := q.qRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	err = q.fRepo.ValidateIsOwner(ctx, foundquestion.Form_id)
	if err != nil {
		return err
	}

	err = q.qRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
