package usecase

import (
	"context"
	"quiz-app/internal/form"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/types"
)

type formUseCase struct {
	formRepo   form.Repo
	ctxUserKey string
}

func NewFormUseCase(formRepo form.Repo, ctxUserKey string) form.UseCase {
	return &formUseCase{
		formRepo:   formRepo,
		ctxUserKey: ctxUserKey,
	}
}

func (f *formUseCase) Create(ctx context.Context, model *models.Form) (*models.Form, error) {
	return f.formRepo.Create(ctx, model)
}

func (f *formUseCase) GetByUserId(ctx context.Context, user_id string, sets types.GetSets) ([]*models.Form, error) {
	return f.formRepo.GetByUserId(ctx, user_id, sets)
}

func (f *formUseCase) Update(ctx context.Context, model *models.Form) (*models.Form, error) {
	err := f.formRepo.ValidateIsOwner(ctx, model.Id)
	if err != nil {
		return nil, err
	}

	currentuser, ok := ctx.Value(f.ctxUserKey).(*models.User)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	// user can not give own form to anyone
	model.User_id = currentuser.Id

	_, err = f.formRepo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (f *formUseCase) Delete(ctx context.Context, id string) error {
	err := f.formRepo.ValidateIsOwner(ctx, id)
	if err != nil {
		return err
	}

	err = f.formRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (f *formUseCase) GetById(ctx context.Context, id string) (*models.Form, error) {
	return f.formRepo.GetById(ctx, id)
}
