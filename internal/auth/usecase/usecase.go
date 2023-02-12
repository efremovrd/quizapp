package usecase

import (
	"context"
	"quiz-app/internal/auth"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/jwter"

	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	authRepo auth.Repo
	jwter    jwter.JWTer
}

func NewAuthUseCase(authRepo auth.Repo, jwter jwter.JWTer) auth.UseCase {
	return &authUseCase{
		authRepo: authRepo,
		jwter:    jwter,
	}
}

func (a *authUseCase) SignUp(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := a.authRepo.GetByLogin(ctx, user.Login)
	if err != errs.ErrContentNotFound {
		if err == nil {
			err = errs.ErrLoginExists
		}
		return nil, err
	}

	nonhashedpswd := user.Password

	pswd, err := bcrypt.GenerateFromPassword([]byte(nonhashedpswd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(pswd)

	createduser, err := a.authRepo.Create(ctx, user)
	if createduser != nil {
		createduser.Password = nonhashedpswd
	}

	return createduser, err
}

func (a *authUseCase) SignIn(ctx context.Context, user *models.User) (*string, error) {
	founduser, err := a.authRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		if err == errs.ErrContentNotFound {
			err = errs.ErrUnauthorized
		}

		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(founduser.Password), []byte(user.Password)) != nil {
		return nil, errs.ErrInvalidPassword
	}

	founduser.Password = user.Password

	return a.jwter.GenerateJWTToken(founduser.Id, founduser.Login)
}

func (a *authUseCase) ParseToken(ctx context.Context, token string) (*models.User, error) {
	id, login, err := a.jwter.ParseToken(token)
	if err != nil {
		return nil, errs.ErrInvalidAccessToken
	}

	return &models.User{
		Id:    *id,
		Login: *login,
	}, nil
}

func (a *authUseCase) GetById(ctx context.Context, id string) (*models.User, error) {
	return a.authRepo.GetById(ctx, id)
}
