package usecase_test

import (
	"context"
	"errors"
	"quiz-app/internal/poolanswer/usecase"
	"quiz-app/models"
	"quiz-app/pkg/types"
	"strconv"
	"testing"

	mocka "quiz-app/internal/answer/mock"
	mockf "quiz-app/internal/form/mock"
	mockpa "quiz-app/internal/poolanswer/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPoolAnswerUseCase_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoA := mocka.NewMockRepo(ctrl)
	mockRepoPA := mockpa.NewMockRepo(ctrl)

	uc := usecase.NewPoolAnswerUseCase(mockRepoPA, mockRepoA, mockRepoF)

	type mockBehavior func(ctx context.Context, pool_answer *models.PoolAnswer, answers []*models.Answer)

	testTable := []struct {
		nameTest        string
		ctx             context.Context
		pool_answer     models.PoolAnswer
		answers         []*models.Answer
		mockBehavior    mockBehavior
		expectedPA      models.PoolAnswer
		expectedAnswers []*models.Answer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			pool_answer: models.PoolAnswer{
				Form_id: "3",
				User_id: "4",
			},
			answers: []*models.Answer{
				{
					Question_id: "7",
					Value:       "ans1",
				},
				{
					Question_id: "8",
					Value:       "ans2",
				},
			},
			mockBehavior: func(ctx context.Context, pool_answer *models.PoolAnswer, answers []*models.Answer) {
				pa := models.PoolAnswer{
					Id:      "10",
					Form_id: pool_answer.Form_id,
					User_id: pool_answer.User_id,
				}
				mockRepoPA.EXPECT().Create(ctx, pool_answer).Return(&pa, nil)
				for i, answer := range answers {
					mockRepoA.EXPECT().Create(ctx, answer).Return(&models.Answer{
						Id:             strconv.Itoa(i),
						Pool_answer_id: pa.Id,
						Question_id:    answer.Question_id,
						Value:          answer.Value,
					}, nil)
				}
			},
			expectedPA: models.PoolAnswer{
				Id:      "10",
				Form_id: "3",
				User_id: "4",
			},
			expectedAnswers: []*models.Answer{
				{
					Question_id:    "7",
					Value:          "ans1",
					Pool_answer_id: "10",
					Id:             "0",
				},
				{
					Question_id:    "8",
					Value:          "ans2",
					Pool_answer_id: "10",
					Id:             "1",
				},
			},
		},
		{
			nameTest: "repoPA_create_error",
			ctx:      context.Background(),
			pool_answer: models.PoolAnswer{
				Form_id: "3",
				User_id: "4",
			},
			answers: []*models.Answer{
				{
					Question_id: "7",
					Value:       "ans1",
				},
				{
					Question_id: "8",
					Value:       "ans2",
				},
			},
			mockBehavior: func(ctx context.Context, pool_answer *models.PoolAnswer, answers []*models.Answer) {
				mockRepoPA.EXPECT().Create(ctx, pool_answer).Return(nil, errors.New("repoPA_create_error"))
			},
		},
		{
			nameTest: "repoA_create_error",
			ctx:      context.Background(),
			pool_answer: models.PoolAnswer{
				Form_id: "3",
				User_id: "4",
			},
			answers: []*models.Answer{
				{
					Question_id: "7",
					Value:       "ans1",
				},
				{
					Question_id: "8",
					Value:       "ans2",
				},
			},
			mockBehavior: func(ctx context.Context, pool_answer *models.PoolAnswer, answers []*models.Answer) {
				pa := models.PoolAnswer{
					Id:      "10",
					Form_id: pool_answer.Form_id,
					User_id: pool_answer.User_id,
				}
				mockRepoPA.EXPECT().Create(ctx, pool_answer).Return(&pa, nil)
				mockRepoA.EXPECT().Create(ctx, answers[0]).Return(nil, errors.New("repoA_create_error"))
				mockRepoPA.EXPECT().Delete(ctx, pa.Id)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.pool_answer, testCase.answers)

			gotpa, gota, err := uc.Create(testCase.ctx, &testCase.pool_answer, testCase.answers)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedPA, *gotpa)
				assert.Equal(t, testCase.expectedAnswers, gota)
			case "repoA_create_error", "repoPA_create_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestPoolAnswerUseCase_GetByFormId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoA := mocka.NewMockRepo(ctrl)
	mockRepoPA := mockpa.NewMockRepo(ctrl)

	uc := usecase.NewPoolAnswerUseCase(mockRepoPA, mockRepoA, mockRepoF)

	type mockBehavior func(ctx context.Context, form_id string, sets types.GetSets)

	testTable := []struct {
		nameTest            string
		ctx                 context.Context
		form_id             string
		sets                types.GetSets
		mockBehavior        mockBehavior
		expectedPoolAnswers []*models.PoolAnswer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			form_id:  "5",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				mockRepoF.EXPECT().ValidateIsOwner(ctx, form_id).Return(nil)
				pas := []*models.PoolAnswer{
					{
						Id:      "10",
						Form_id: form_id,
						User_id: "32",
					},
					{
						Id:      "11",
						Form_id: form_id,
						User_id: "33",
					},
				}
				mockRepoPA.EXPECT().GetByFormId(ctx, form_id, sets).Return(pas, nil)
			},
			expectedPoolAnswers: []*models.PoolAnswer{
				{
					Id:      "10",
					Form_id: "5",
					User_id: "32",
				},
				{
					Id:      "11",
					Form_id: "5",
					User_id: "33",
				},
			},
		},
		{
			nameTest: "user_is_not_an_owner",
			ctx:      context.Background(),
			form_id:  "5",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				mockRepoF.EXPECT().ValidateIsOwner(ctx, form_id).Return(errors.New("user_is_not_an_owner"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.form_id, testCase.sets)

			gotpas, err := uc.GetByFormId(testCase.ctx, testCase.form_id, testCase.sets)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedPoolAnswers, gotpas)
			case "user_is_not_an_owner":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestPoolAnswerUseCase_GetId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoA := mocka.NewMockRepo(ctrl)
	mockRepoPA := mockpa.NewMockRepo(ctrl)

	uc := usecase.NewPoolAnswerUseCase(mockRepoPA, mockRepoA, mockRepoF)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest           string
		ctx                context.Context
		id                 string
		mockBehavior       mockBehavior
		expectedPoolAnswer models.PoolAnswer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "5",
			mockBehavior: func(ctx context.Context, id string) {
				pa := models.PoolAnswer{
					Id:      id,
					Form_id: "10",
					User_id: "32",
				}
				mockRepoPA.EXPECT().GetById(ctx, id).Return(&pa, nil)
			},
			expectedPoolAnswer: models.PoolAnswer{
				Id:      "5",
				Form_id: "10",
				User_id: "32",
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
				assert.Equal(t, testCase.expectedPoolAnswer, *got)
			case "user_is_not_an_owner":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
