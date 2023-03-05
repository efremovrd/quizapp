package usecase_test

import (
	"context"
	"errors"
	"quizapp/internal/form/mock"
	"quizapp/internal/form/usecase"
	"quizapp/models"
	"quizapp/pkg/errs"
	"quizapp/pkg/types"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFormUseCase_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepo(ctrl)

	ctxUserKey := "ctxuserkey"

	uc := usecase.NewFormUseCase(mockRepo, ctxUserKey)

	type mockBehavior func(ctx context.Context, model *models.Form)

	testTable := []struct {
		nameTest      string
		ctx           context.Context
		model         models.Form
		mockBehavior  mockBehavior
		expectedModel models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			model: models.Form{
				User_id:     "5",
				Title:       "title",
				Description: "desc",
			},
			mockBehavior: func(ctx context.Context, model *models.Form) {
				mockRepo.EXPECT().Create(ctx, model).Return(&models.Form{
					Id:          "1",
					User_id:     model.User_id,
					Title:       model.Title,
					Description: model.Description,
				}, nil)
			},
			expectedModel: models.Form{
				Id:          "1",
				User_id:     "5",
				Title:       "title",
				Description: "desc",
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.model)

			got, err := uc.Create(testCase.ctx, &testCase.model)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedModel, *got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormUseCase_GetByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepo(ctrl)

	ctxUserKey := "ctxuserkey"

	uc := usecase.NewFormUseCase(mockRepo, ctxUserKey)

	type mockBehavior func(ctx context.Context, user_id string, sets types.GetSets)

	testTable := []struct {
		nameTest       string
		ctx            context.Context
		user_id        string
		sets           types.GetSets
		mockBehavior   mockBehavior
		expectedModels []*models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			user_id:  "5",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, user_id string, sets types.GetSets) {
				mockRepo.EXPECT().GetByUserId(ctx, user_id, sets).Return([]*models.Form{
					{
						Id:          "1",
						User_id:     user_id,
						Title:       "title1",
						Description: "desc1",
					},
					{
						Id:          "2",
						User_id:     user_id,
						Title:       "title2",
						Description: "desc2",
					},
				}, nil)
			},
			expectedModels: []*models.Form{
				{
					Id:          "1",
					User_id:     "5",
					Title:       "title1",
					Description: "desc1",
				},
				{
					Id:          "2",
					User_id:     "5",
					Title:       "title2",
					Description: "desc2",
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.user_id, testCase.sets)

			got, err := uc.GetByUserId(testCase.ctx, testCase.user_id, testCase.sets)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedModels, got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormUseCase_GetById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepo(ctrl)

	ctxUserKey := "ctxuserkey"

	uc := usecase.NewFormUseCase(mockRepo, ctxUserKey)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest      string
		ctx           context.Context
		id            string
		mockBehavior  mockBehavior
		expectedModel models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "1",
			mockBehavior: func(ctx context.Context, id string) {
				mockRepo.EXPECT().GetById(ctx, id).Return(&models.Form{
					Id:          id,
					User_id:     "5",
					Title:       "title1",
					Description: "desc1",
				}, nil)
			},
			expectedModel: models.Form{
				Id:          "1",
				User_id:     "5",
				Title:       "title1",
				Description: "desc1",
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
				assert.Equal(t, testCase.expectedModel, *got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormUseCase_Delete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepo(ctrl)

	ctxUserKey := "ctxuserkey"

	uc := usecase.NewFormUseCase(mockRepo, ctxUserKey)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		id           string
		mockBehavior mockBehavior
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "1",
			mockBehavior: func(ctx context.Context, id string) {
				mockRepo.EXPECT().ValidateIsOwner(ctx, id).Return(nil)
				mockRepo.EXPECT().Delete(ctx, id).Return(nil)
			},
		},
		{
			nameTest: "repo_delete_error",
			ctx:      context.Background(),
			id:       "1",
			mockBehavior: func(ctx context.Context, id string) {
				mockRepo.EXPECT().ValidateIsOwner(ctx, id).Return(nil)
				mockRepo.EXPECT().Delete(ctx, id).Return(errors.New("repo_delete_error"))
			},
		},
		{
			nameTest: "user_is_not_an_owner",
			ctx:      context.Background(),
			id:       "1",
			mockBehavior: func(ctx context.Context, id string) {
				mockRepo.EXPECT().ValidateIsOwner(ctx, id).Return(errs.ErrForbidden)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.id)

			err := uc.Delete(testCase.ctx, testCase.id)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
			case "repo_delete_error":
				assert.NotEqual(t, nil, err)
			case "user_is_not_an_owner":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormUseCase_Update(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepo(ctrl)

	ctxUserKey := "ctxuserkey"

	uc := usecase.NewFormUseCase(mockRepo, ctxUserKey)

	type mockBehavior func(ctx context.Context, model *models.Form)

	testTable := []struct {
		nameTest      string
		ctx           context.Context
		model         models.Form
		mockBehavior  mockBehavior
		expectedModel models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.WithValue(context.Background(), ctxUserKey, &models.User{Id: "5"}),
			model: models.Form{
				Id:          "1",
				User_id:     "5",
				Title:       "title",
				Description: "desc",
			},
			mockBehavior: func(ctx context.Context, model *models.Form) {
				mockRepo.EXPECT().ValidateIsOwner(ctx, model.Id).Return(nil)
				mockRepo.EXPECT().Update(ctx, model).Return(model, nil)
			},
			expectedModel: models.Form{
				Id:          "1",
				User_id:     "5",
				Title:       "title",
				Description: "desc",
			},
		},
		{
			nameTest: "user_is_not_an_owner",
			ctx:      context.Background(),
			model: models.Form{
				Id:          "1",
				User_id:     "5",
				Title:       "title",
				Description: "desc",
			},
			mockBehavior: func(ctx context.Context, model *models.Form) {
				mockRepo.EXPECT().ValidateIsOwner(ctx, model.Id).Return(errors.New("user_is_not_an_owner"))
			},
		},
		{
			nameTest: "repo_update_error",
			ctx:      context.WithValue(context.Background(), ctxUserKey, &models.User{Id: "5"}),
			model: models.Form{
				Id:          "1",
				User_id:     "5",
				Title:       "title",
				Description: "desc",
			},
			mockBehavior: func(ctx context.Context, model *models.Form) {
				mockRepo.EXPECT().ValidateIsOwner(ctx, model.Id).Return(nil)
				mockRepo.EXPECT().Update(ctx, model).Return(nil, errors.New("repo_update_error"))
			},
		},
		{
			nameTest: "unauthorized",
			ctx:      context.Background(),
			model: models.Form{
				Id:          "1",
				User_id:     "5",
				Title:       "title",
				Description: "desc",
			},
			mockBehavior: func(ctx context.Context, model *models.Form) {
				mockRepo.EXPECT().ValidateIsOwner(ctx, model.Id).Return(nil)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.model)

			got, err := uc.Update(testCase.ctx, &testCase.model)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedModel, *got)
			case "user_is_not_an_owner", "repo_update_error":
				assert.NotEqual(t, nil, err)
			case "unauthorized":
				assert.Equal(t, errs.ErrUnauthorized, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
