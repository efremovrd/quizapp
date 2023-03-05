package usecase_test

import (
	"context"
	"errors"
	mockf "quizapp/internal/form/mock"
	mockq "quizapp/internal/question/mock"
	"quizapp/internal/question/usecase"
	"quizapp/models"
	"quizapp/pkg/types"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestQuestionUseCase_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoQ := mockq.NewMockRepo(ctrl)

	uc := usecase.NewQuestionUseCase(mockRepoQ, mockRepoF)

	type mockBehavior func(ctx context.Context, model *models.Question)

	testTable := []struct {
		nameTest      string
		ctx           context.Context
		model         models.Question
		mockBehavior  mockBehavior
		expectedModel models.Question
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			model: models.Question{
				Form_id: "5",
				Header:  "header",
			},
			mockBehavior: func(ctx context.Context, model *models.Question) {
				mockRepoF.EXPECT().ValidateIsOwner(ctx, model.Form_id).Return(nil)
				mockRepoQ.EXPECT().Create(ctx, model).Return(&models.Question{
					Id:      "1",
					Form_id: model.Form_id,
					Header:  model.Header,
				}, nil)
			},
			expectedModel: models.Question{
				Id:      "1",
				Form_id: "5",
				Header:  "header",
			},
		},
		{
			nameTest: "user_is_not_an_owner",
			ctx:      context.Background(),
			model: models.Question{
				Form_id: "5",
				Header:  "header",
			},
			mockBehavior: func(ctx context.Context, model *models.Question) {
				mockRepoF.EXPECT().ValidateIsOwner(ctx, model.Form_id).Return(errors.New("user_is_not_an_owner"))
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
			case "user_is_not_an_owner":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestQuestionUseCase_GetByFormId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoQ := mockq.NewMockRepo(ctrl)

	uc := usecase.NewQuestionUseCase(mockRepoQ, mockRepoF)

	type mockBehavior func(ctx context.Context, form_id string, sets types.GetSets)

	testTable := []struct {
		nameTest       string
		ctx            context.Context
		form_id        string
		sets           types.GetSets
		mockBehavior   mockBehavior
		expectedModels []*models.Question
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			form_id:  "5",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				mockRepoQ.EXPECT().GetByFormId(ctx, form_id, sets).Return([]*models.Question{
					{
						Id:      "1",
						Form_id: form_id,
						Header:  "header1",
					},
					{
						Id:      "2",
						Form_id: form_id,
						Header:  "header2",
					},
				}, nil)
			},
			expectedModels: []*models.Question{
				{
					Id:      "1",
					Form_id: "5",
					Header:  "header1",
				},
				{
					Id:      "2",
					Form_id: "5",
					Header:  "header2",
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.form_id, testCase.sets)

			got, err := uc.GetByFormId(testCase.ctx, testCase.form_id, testCase.sets)

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

func TestQuestionUseCase_Update(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoQ := mockq.NewMockRepo(ctrl)

	uc := usecase.NewQuestionUseCase(mockRepoQ, mockRepoF)

	type mockBehavior func(ctx context.Context, model *models.Question)

	testTable := []struct {
		nameTest      string
		ctx           context.Context
		model         models.Question
		mockBehavior  mockBehavior
		expectedModel models.Question
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			model: models.Question{
				Id:      "1",
				Form_id: "5",
				Header:  "header",
			},
			mockBehavior: func(ctx context.Context, model *models.Question) {
				mockRepoF.EXPECT().ValidateIsOwner(ctx, model.Form_id).Return(nil)
				mockRepoQ.EXPECT().Update(ctx, model).Return(nil, nil)
			},
			expectedModel: models.Question{
				Id:      "1",
				Form_id: "5",
				Header:  "header",
			},
		},
		{
			nameTest: "user_is_not_an_owner",
			ctx:      context.Background(),
			model: models.Question{
				Id:      "1",
				Form_id: "5",
				Header:  "header",
			},
			mockBehavior: func(ctx context.Context, model *models.Question) {
				mockRepoF.EXPECT().ValidateIsOwner(ctx, model.Form_id).Return(errors.New("user_is_not_an_owner"))
			},
		},
		{
			nameTest: "repo_update_error",
			ctx:      context.Background(),
			model: models.Question{
				Id:      "1",
				Form_id: "5",
				Header:  "header",
			},
			mockBehavior: func(ctx context.Context, model *models.Question) {
				mockRepoF.EXPECT().ValidateIsOwner(ctx, model.Form_id).Return(nil)
				mockRepoQ.EXPECT().Update(ctx, model).Return(nil, errors.New("repo_update_error"))
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
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestQuestionUseCase_Delete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoQ := mockq.NewMockRepo(ctrl)

	uc := usecase.NewQuestionUseCase(mockRepoQ, mockRepoF)

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
				formid := "5"
				mockRepoQ.EXPECT().GetById(ctx, id).Return(&models.Question{
					Id:      id,
					Form_id: formid,
					Header:  "header",
				}, nil)
				mockRepoF.EXPECT().ValidateIsOwner(ctx, formid).Return(nil)
				mockRepoQ.EXPECT().Delete(ctx, id).Return(nil)
			},
		},
		{
			nameTest: "user_is_not_an_owner",
			ctx:      context.Background(),
			id:       "1",
			mockBehavior: func(ctx context.Context, id string) {
				formid := "5"
				mockRepoQ.EXPECT().GetById(ctx, id).Return(&models.Question{
					Id:      id,
					Form_id: formid,
					Header:  "header",
				}, nil)
				mockRepoF.EXPECT().ValidateIsOwner(ctx, formid).Return(errors.New("user_is_not_an_owner"))
			},
		},
		{
			nameTest: "repo_delete_error",
			ctx:      context.Background(),
			id:       "1",
			mockBehavior: func(ctx context.Context, id string) {
				formid := "5"
				mockRepoQ.EXPECT().GetById(ctx, id).Return(&models.Question{
					Id:      id,
					Form_id: formid,
					Header:  "header",
				}, nil)
				mockRepoF.EXPECT().ValidateIsOwner(ctx, formid).Return(nil)
				mockRepoQ.EXPECT().Delete(ctx, id).Return(errors.New("repo_delete_error"))
			},
		},
		{
			nameTest: "no_question_to_delete",
			ctx:      context.Background(),
			id:       "2",
			mockBehavior: func(ctx context.Context, id string) {
				mockRepoQ.EXPECT().GetById(ctx, id).Return(nil, errors.New("no_question_to_delete"))
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
			case "user_is_not_an_owner", "repo_delete_error", "no_question_to_delete":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
