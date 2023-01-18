package usecase_test

import (
	"context"
	"errors"
	"quiz-app/internal/auth/usecase"
	"quiz-app/models"
	"testing"

	mockauth "quiz-app/internal/auth/mock"
	"quiz-app/pkg/errs"
	mockjwt "quiz-app/pkg/jwter/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUseCase_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoAuth := mockauth.NewMockRepo(ctrl)
	mockjwter := mockjwt.NewMockJWTer(ctrl)

	uc := usecase.NewAuthUseCase(mockRepoAuth, mockjwter)

	type mockBehavior func(ctx context.Context, user *models.User)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		user         models.User
		mockBehavior mockBehavior
		expectedUser models.User
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				createduser := models.User{
					Id:       "5",
					Login:    user.Login,
					Password: user.Password,
				}
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(nil, errs.ErrContentNotFound)
				mockRepoAuth.EXPECT().Create(ctx, user).Return(&createduser, nil)
			},
			expectedUser: models.User{
				Id:       "5",
				Login:    "login",
				Password: "password",
			},
		},
		{
			nameTest: "repoAuth_create_error",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(nil, errs.ErrContentNotFound)
				mockRepoAuth.EXPECT().Create(ctx, user).Return(nil, errors.New("repoAuth_create_error"))
			},
		},
		{
			nameTest: "login_exists",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(nil, nil)
			},
		},
		{
			nameTest: "repoAuth_getbylogin_error",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(nil, errors.New("repoAuth_getbylogin_error"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.user)

			got, err := uc.SignUp(testCase.ctx, &testCase.user)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedUser, *got)
			case "repoAuth_create_error", "repoAuth_getbylogin_error":
				assert.NotEqual(t, nil, err)
			case "login_exists":
				assert.Equal(t, errs.ErrLoginExists, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestAuthUseCase_SignIn(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoAuth := mockauth.NewMockRepo(ctrl)
	mockjwter := mockjwt.NewMockJWTer(ctrl)

	uc := usecase.NewAuthUseCase(mockRepoAuth, mockjwter)

	type mockBehavior func(ctx context.Context, user *models.User)

	testTable := []struct {
		nameTest      string
		ctx           context.Context
		user          models.User
		mockBehavior  mockBehavior
		expectedToken string
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				pswd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				founduser := models.User{
					Id:       "5",
					Login:    user.Login,
					Password: string(pswd),
				}
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(&founduser, nil)
				token := "token"
				mockjwter.EXPECT().GenerateJWTToken(&models.User{
					Id:       founduser.Id,
					Login:    founduser.Login,
					Password: user.Password,
				}).Return(&token, nil)
			},
			expectedToken: "token",
		},
		{
			nameTest: "no_such_user",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(nil, errs.ErrContentNotFound)
			},
		},
		{
			nameTest: "repoAuth_getbylogin_error",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(nil, errors.New("repoAuth_getbylogin_error"))
			},
		},
		{
			nameTest: "wrong_password",
			ctx:      context.Background(),
			user: models.User{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				founduser := models.User{
					Id:       "5",
					Login:    user.Login,
					Password: "anotherpassword",
				}
				mockRepoAuth.EXPECT().GetByLogin(ctx, user.Login).Return(&founduser, nil)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.user)

			got, err := uc.SignIn(testCase.ctx, &testCase.user)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedToken, *got)
			case "repoAuth_getbylogin_error":
				assert.NotEqual(t, nil, err)
			case "wrong_password":
				assert.Equal(t, errs.ErrInvalidPassword, err)
			case "no_such_user":
				assert.Equal(t, errs.ErrUnauthorized, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestAuthUseCase_ParseToken(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoAuth := mockauth.NewMockRepo(ctrl)
	mockjwter := mockjwt.NewMockJWTer(ctrl)

	uc := usecase.NewAuthUseCase(mockRepoAuth, mockjwter)

	type mockBehavior func(token string)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		token        string
		mockBehavior mockBehavior
		expectedUser models.User
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			token:    "token",
			mockBehavior: func(token string) {
				user := models.User{
					Id:       "5",
					Login:    "login",
					Password: "password",
				}
				mockjwter.EXPECT().ParseToken(token).Return(&user, nil)
			},
			expectedUser: models.User{
				Id:       "5",
				Login:    "login",
				Password: "password",
			},
		},
		{
			nameTest: "invalid_token",
			ctx:      context.Background(),
			token:    "invalidtoken",
			mockBehavior: func(token string) {
				mockjwter.EXPECT().ParseToken(token).Return(nil, errors.New("invalid_token"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.token)

			got, err := uc.ParseToken(testCase.ctx, testCase.token)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedUser, *got)
			case "invalid_token":
				assert.Equal(t, errs.ErrInvalidAccessToken, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestAuthUseCase_GetById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoAuth := mockauth.NewMockRepo(ctrl)
	mockjwter := mockjwt.NewMockJWTer(ctrl)

	uc := usecase.NewAuthUseCase(mockRepoAuth, mockjwter)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		id           string
		mockBehavior mockBehavior
		expectedUser models.User
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "2",
			mockBehavior: func(ctx context.Context, id string) {
				user := models.User{
					Id:       id,
					Login:    "login",
					Password: "password",
				}
				mockRepoAuth.EXPECT().GetById(ctx, id).Return(&user, nil)
			},
			expectedUser: models.User{
				Id:       "2",
				Login:    "login",
				Password: "password",
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.id)

			got, err := uc.GetById(testCase.ctx, testCase.id)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedUser, *got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
