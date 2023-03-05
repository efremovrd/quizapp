package usecase_test

import (
	"context"
	"errors"
	"quizapp/internal/answer/usecase"
	"quizapp/models"
	"quizapp/pkg/types"
	"testing"

	mocka "quiz-app/internal/answer/mock"
	mockf "quiz-app/internal/form/mock"
	mockpa "quiz-app/internal/poolanswer/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAnswerUseCase_GetByPoolAnswerId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoF := mockf.NewMockRepo(ctrl)
	mockRepoA := mocka.NewMockRepo(ctrl)
	mockRepoPA := mockpa.NewMockRepo(ctrl)

	uc := usecase.NewAnswerUseCase(mockRepoA, mockRepoF, mockRepoPA)

	type mockBehavior func(ctx context.Context, pool_answer_id string, sets types.GetSets)

	testTable := []struct {
		nameTest        string
		ctx             context.Context
		pool_answer_id  string
		sets            types.GetSets
		mockBehavior    mockBehavior
		expectedAnswers []*models.Answer
	}{
		{
			nameTest:       "ok",
			ctx:            context.Background(),
			pool_answer_id: "5",
			sets:           types.GetSets{},
			mockBehavior: func(ctx context.Context, pool_answer_id string, sets types.GetSets) {
				foundpa := models.PoolAnswer{
					Id:      pool_answer_id,
					Form_id: "45",
					User_id: "46",
				}
				mockRepoPA.EXPECT().GetById(ctx, pool_answer_id).Return(&foundpa, nil)
				mockRepoF.EXPECT().ValidateIsOwner(ctx, foundpa.Form_id).Return(nil)
				mockRepoA.EXPECT().GetByPoolAnswerId(ctx, pool_answer_id, sets).Return([]*models.Answer{
					{
						Question_id:    "7",
						Value:          "ans1",
						Pool_answer_id: pool_answer_id,
						Id:             "0",
					},
					{
						Question_id:    "8",
						Value:          "ans2",
						Pool_answer_id: pool_answer_id,
						Id:             "1",
					},
				}, nil)
			},
			expectedAnswers: []*models.Answer{
				{
					Question_id:    "7",
					Value:          "ans1",
					Pool_answer_id: "5",
					Id:             "0",
				},
				{
					Question_id:    "8",
					Value:          "ans2",
					Pool_answer_id: "5",
					Id:             "1",
				},
			},
		},
		{
			nameTest:       "paRepo_getbyid_error",
			ctx:            context.Background(),
			pool_answer_id: "5",
			sets:           types.GetSets{},
			mockBehavior: func(ctx context.Context, pool_answer_id string, sets types.GetSets) {
				mockRepoPA.EXPECT().GetById(ctx, pool_answer_id).Return(nil, errors.New("paRepo_getbyid_error"))
			},
		},
		{
			nameTest:       "user_not_an_owner",
			ctx:            context.Background(),
			pool_answer_id: "5",
			sets:           types.GetSets{},
			mockBehavior: func(ctx context.Context, pool_answer_id string, sets types.GetSets) {
				foundpa := models.PoolAnswer{
					Id:      pool_answer_id,
					Form_id: "45",
					User_id: "46",
				}
				mockRepoPA.EXPECT().GetById(ctx, pool_answer_id).Return(&foundpa, nil)
				mockRepoF.EXPECT().ValidateIsOwner(ctx, foundpa.Form_id).Return(errors.New("user_not_an_owner"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.pool_answer_id, testCase.sets)

			got, err := uc.GetByPoolAnswerId(testCase.ctx, testCase.pool_answer_id, testCase.sets)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedAnswers, got)
			case "paRepo_getbyid_error", "user_not_an_owner":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
