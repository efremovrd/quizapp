package http

import (
	"net/http"
	"quiz-app/internal/question"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/types"

	"github.com/gin-gonic/gin"
)

type questionCreatRequest struct {
	Header string `json:"header" binding:"required"`
}

type questionResponse struct {
	Id      string `json:"id"`
	Form_id string `json:"form_id"`
	Header  string `json:"header"`
}

type questionGetByFormIdResponse struct {
	Questions []*questionResponse `json:"questions"`
}

type questionHandlers struct {
	questionUC question.UseCase
	ctxUserKey string
}

func NewQuestionHandlers(questionUC question.UseCase, ctxUserKey string) question.Handlers {
	return &questionHandlers{
		questionUC: questionUC,
		ctxUserKey: ctxUserKey,
	}
}

// Create godoc
// @Summary Create question
// @Description Create new question with header for current form
// @Tags Questions
// @Security JWTToken
// @Param data body questionCreatRequest true "question header"
// @Param id path string true "current form id"
// @Success 201 {object} questionResponse
// @Failure 204   "No such form"
// @Failure 400   "Invalid json"
// @Failure 401   "Unauthorized"
// @Failure 403   "Permission denied"
// @Failure 500   "Other err"
// @Router /forms/{id}/questions [post]
func (h *questionHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(questionCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := &models.Question{
			Form_id: c.Param("formid"),
			Header:  request.Header,
		}

		createdquestion, err := h.questionUC.Create(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusCreated, questionBLToResponse(createdquestion))
	}
}

// Delete godoc
// @Summary Delete question
// @Description Delete question by id
// @Tags Questions
// @Security JWTToken
// @Param id path string true "question id"
// @Success 200   "Deleted"
// @Failure 204   "No such question"
// @Failure 400   "Invalid question id"
// @Failure 401   "Unauthorized"
// @Failure 403   "User is not the form owner or permission denied"
// @Failure 500   "Other err"
// @Router /forms/{formid}/questions/{questionid} [delete]
func (h *questionHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h.questionUC.Delete(c, c.Param("questionid"))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Status(http.StatusOK)
	}
}

// GetByFormId godoc
// @Summary Get questions
// @Description Get questions by form
// @Tags Questions
// @Security JWTToken
// @Param limit query int true "limit" minimum(1)
// @Param offset query int true "offset" minimum(0)
// @Param formid path string true "form id"
// @Success 200 {object} questionGetByFormIdResponse "Found"
// @Failure 204 {object} questionGetByFormIdResponse "No questions by form"
// @Failure 400   "Invalid params"
// @Failure 401   "Unauthorized"
// @Failure 500   "Other err"
// @Router /forms/{formid}/questions [get]
func (h *questionHandlers) GetByFormId() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset, ok := types.ValidateGetSets(c.Query("limit"), c.Query("offset"))
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundquestions, err := h.questionUC.GetByFormId(c, c.Param("formid"), types.GetSets{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		status := http.StatusOK

		if len(foundquestions) == 0 {
			status = http.StatusNoContent
		}

		c.JSON(status, &questionGetByFormIdResponse{
			Questions: questionsBLToResponse(foundquestions),
		})
	}
}

// Update godoc
// @Summary Update question
// @Description Update question with header
// @Tags Questions
// @Security JWTToken
// @Param formid path string true "form id"
// @Param questionid path string true "question id"
// @Param new body questionCreatRequest true "new header"
// @Success 200 {object} questionResponse "Updated"
// @Failure 204   "No such question"
// @Failure 400   "Invalid question id"
// @Failure 401   "Unauthorized"
// @Failure 403   "Permission denied"
// @Failure 500   "Other err"
// @Router /forms/{formid}/questions/{questionid} [put]
func (h *questionHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(questionCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		updatedquestion, err := h.questionUC.Update(c, &models.Question{
			Id:      c.Param("questionid"),
			Form_id: c.Param("formid"),
			Header:  request.Header,
		})
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, questionBLToResponse(updatedquestion))
	}
}

func questionBLToResponse(modelBL *models.Question) *questionResponse {
	return &questionResponse{
		Id:      modelBL.Id,
		Form_id: modelBL.Form_id,
		Header:  modelBL.Header,
	}
}

func questionsBLToResponse(questions []*models.Question) []*questionResponse {
	if questions == nil {
		return nil
	}

	res := make([]*questionResponse, len(questions))

	for i, q := range questions {
		res[i] = questionBLToResponse(q)
	}

	return res
}
