package http

import (
	"net/http"
	"quiz-app/internal/answer"
	"quiz-app/internal/form"
	"quiz-app/internal/poolanswer"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/types"

	"github.com/gin-gonic/gin"
)

type answerResponse struct {
	Id             string `json:"id"`
	Question_id    string `json:"question_id"`
	Pool_answer_id string `json:"pool_answer_id"`
	Value          string `json:"value"`
}

type poolAnswerResponse struct {
	Id      string `json:"id"`
	User_id string `json:"user_id"`
	Form_id string `json:"form_id"`
}

type poolsAnswerResponse struct {
	Pools_answer []*poolAnswerResponse `json:"pools_answer"`
}

type answersResponse struct {
	Pool_answer *poolAnswerResponse `json:"pool_answer"`
	Answers     []*answerResponse   `json:"answers"`
}

type answerRequest struct {
	Question_id string `json:"question_id" binding:"required"`
	Value       string `json:"value" binding:"required"`
}

type poolAnswerCreatRequest struct {
	Answers []*answerRequest `json:"answers" binding:"required"`
}

type poolAnswerCreatResponse struct {
	Pool_answer *poolAnswerResponse `json:"pool_answer"`
	Answers     []*answerResponse   `json:"answers"`
}

type answersHandlers struct {
	paUC       poolanswer.UseCase
	aUC        answer.UseCase
	fUC        form.UseCase
	ctxUserKey string
}

func NewAnswersHandlers(paUC poolanswer.UseCase, aUC answer.UseCase, fUC form.UseCase, ctxUserKey string) poolanswer.Handlers {
	return &answersHandlers{
		paUC:       paUC,
		aUC:        aUC,
		fUC:        fUC,
		ctxUserKey: ctxUserKey,
	}
}

// Create godoc
// @Summary Create answers
// @Description Create answers: pool answer by form id, answers with question id and value
// @Tags Answers
// @Security JWTToken
// @Param data body poolAnswerCreatRequest true "answers"
// @Param formid path string true "form id"
// @Success 201 {object} poolAnswerCreatResponse
// @Failure 400   "Invalid params"
// @Failure 401   "Unauthorized"
// @Failure 403   "Permission denied"
// @Failure 500   "Other err"
// @Router /forms/{formid}/poolsanswer [post]
func (h *answersHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentuser, ok := c.Value(h.ctxUserKey).(*models.User)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		answersDTO := new(poolAnswerCreatRequest)

		err := c.ShouldBindJSON(answersDTO)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		poolanswer := &models.PoolAnswer{
			User_id: currentuser.Id,
			Form_id: c.Param("formid"),
		}

		createdpa, createdanswers, err := h.paUC.Create(c, poolanswer, answersDTOToBL(answersDTO.Answers))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusCreated, &poolAnswerCreatResponse{
			Pool_answer: poolAnswerBLToDTO(createdpa),
			Answers:     answersBLToDTO(createdanswers),
		})
	}
}

// GetByFormId godoc
// @Summary Get answers
// @Description Get pool answer by form
// @Tags Answers
// @Security JWTToken
// @Param limit query int true "limit" minimum(1)
// @Param offset query int true "offset" minimum(0)
// @Param formid path string true "form id"
// @Success 200 {object} poolsAnswerResponse "Found"
// @Failure 204 {object} poolsAnswerResponse "No such form or pools answer"
// @Failure 400   "Invalid params"
// @Failure 401   "Unauthorized"
// @Failure 403   "User is not the form owner"
// @Failure 500   "Other err"
// @Router /forms/{formid}/poolsanswer [get]
func (h *answersHandlers) GetByFormId() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset, ok := types.ValidateGetSets(c.Query("limit"), c.Query("offset"))
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		pools_answer, err := h.paUC.GetByFormId(c, c.Param("formid"), types.GetSets{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		status := http.StatusOK

		if len(pools_answer) == 0 {
			status = http.StatusNoContent
		}

		c.JSON(status, &poolsAnswerResponse{
			Pools_answer: poolsanswerBLToDTO(pools_answer),
		})
	}
}

// GetByPoolAnswerId godoc
// @Summary Get answers
// @Description Get answers by pool answer id
// @Tags Answers
// @Security JWTToken
// @Param poolanswerid path string true "pool answer id"
// @Param limit query int true "limit" minimum(1)
// @Param offset query int true "offset" minimum(0)
// @Param formid path string true "form id"
// @Success 200 {object} answersResponse "Found"
// @Failure 204 {object} answersResponse "No such answers"
// @Failure 400   "Invalid params"
// @Failure 401   "Unauthorized"
// @Failure 403   "User is not the form owner"
// @Failure 500   "Other err"
// @Router /forms/{formid}/poolsanswer/{poolanswerid} [get]
func (h *answersHandlers) GetByPoolAnswerId() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset, ok := types.ValidateGetSets(c.Query("limit"), c.Query("offset"))
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		pa, err := h.paUC.GetById(c, c.Param("poolanswerid"))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		if pa.Form_id != c.Param("formid") {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		foundanswers, err := h.aUC.GetByPoolAnswerId(c, c.Param("poolanswerid"), types.GetSets{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		status := http.StatusOK

		if len(foundanswers) == 0 {
			status = http.StatusNoContent
		}

		c.JSON(status, &answersResponse{
			Pool_answer: poolAnswerBLToDTO(pa),
			Answers:     answersBLToDTO(foundanswers),
		})
	}
}

func answerBLToDTO(answerBL *models.Answer) *answerResponse {
	return &answerResponse{
		Id:             answerBL.Id,
		Question_id:    answerBL.Question_id,
		Pool_answer_id: answerBL.Pool_answer_id,
		Value:          answerBL.Value,
	}
}

func answerDTOToBL(answerDTO *answerRequest) *models.Answer {
	return &models.Answer{
		Question_id: answerDTO.Question_id,
		Value:       answerDTO.Value,
	}
}

func answersDTOToBL(answersDTO []*answerRequest) []*models.Answer {
	res := make([]*models.Answer, len(answersDTO))

	for i, a := range answersDTO {
		res[i] = answerDTOToBL(a)
	}

	return res
}

func answersBLToDTO(answers []*models.Answer) []*answerResponse {
	res := make([]*answerResponse, len(answers))

	for i, a := range answers {
		res[i] = answerBLToDTO(a)
	}

	return res
}

func poolAnswerBLToDTO(paBL *models.PoolAnswer) *poolAnswerResponse {
	return &poolAnswerResponse{
		Id:      paBL.Id,
		Form_id: paBL.Form_id,
		User_id: paBL.User_id,
	}
}

func poolsanswerBLToDTO(paBL []*models.PoolAnswer) []*poolAnswerResponse {
	res := make([]*poolAnswerResponse, len(paBL))

	for i, a := range paBL {
		res[i] = poolAnswerBLToDTO(a)
	}

	return res
}
