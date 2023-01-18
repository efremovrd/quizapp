package http

import (
	"net/http"
	"quiz-app/internal/form"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/types"

	"github.com/gin-gonic/gin"
)

type formCreatRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type formUpdRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type formResponse struct {
	Id          string `json:"id"`
	User_id     string `json:"user_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type formGetByUserIdResponse struct {
	Forms []*formResponse `json:"forms"`
}

type formHandlers struct {
	formUC     form.UseCase
	ctxUserKey string
}

func NewFormHandlers(formUC form.UseCase, ctxUserKey string) form.Handlers {
	return &formHandlers{
		formUC:     formUC,
		ctxUserKey: ctxUserKey,
	}
}

// Create godoc
// @Summary Create form
// @Description Create new form with title and description
// @Tags Forms
// @Security JWTToken
// @Accept json
// @Param data body formCreatRequest true "form title and description"
// @Success 201 {object} formResponse
// @Failure 400   "Invalid json"
// @Failure 401   "Unauthorized"
// @Failure 403   "Permission denied"
// @Failure 500   "Other err"
// @Router /forms [post]
func (h *formHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentuser, ok := c.Value(h.ctxUserKey).(*models.User)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		request := new(formCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := formCreatRequestToBL(request)
		modelBL.User_id = currentuser.Id

		createdform, err := h.formUC.Create(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusCreated, formBLToResponse(createdform))
	}
}

// Delete godoc
// @Summary Delete form
// @Description Delete form by id
// @Tags Forms
// @Security JWTToken
// @Param formid path string true "form id"
// @Success 200   "Deleted"
// @Failure 204   "No such form"
// @Failure 400   "Invalid id"
// @Failure 401   "Unauthorized"
// @Failure 403   "Permission denied"
// @Failure 500   "Other err"
// @Router /forms/{formid} [delete]
func (h *formHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h.formUC.Delete(c, c.Param("formid"))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Status(http.StatusOK)
	}
}

// GetByUser godoc
// @Summary Get forms
// @Description Get forms owned by current user
// @Tags Forms
// @Security JWTToken
// @Param limit query int true "limit" minimum(1)
// @Param offset query int true "offset" minimum(0)
// @Success 200 {object} formGetByUserIdResponse "Found"
// @Failure 204 {object} formGetByUserIdResponse "No forms owned by current user"
// @Failure 400   "Invalid limit and/or offset"
// @Failure 401   "Unauthorized"
// @Failure 500   "Other err"
// @Router /forms [get]
func (h *formHandlers) GetByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentuser, ok := c.Value(h.ctxUserKey).(*models.User)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		limit, offset, ok := types.ValidateGetSets(c.Query("limit"), c.Query("offset"))
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundforms, err := h.formUC.GetByUserId(c, currentuser.Id, types.GetSets{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		status := http.StatusOK

		if len(foundforms) == 0 {
			status = http.StatusNoContent
		}

		c.JSON(status, &formGetByUserIdResponse{
			Forms: formsBLToResponse(foundforms),
		})
	}
}

// Update godoc
// @Summary Update form
// @Description Update form with title and/or description
// @Tags Forms
// @Security JWTToken
// @Param formid path string true "form id"
// @Param new body formUpdRequest true "new title and/or description"
// @Success 200 {object} formResponse "Updated"
// @Failure 204   "No such form"
// @Failure 400   "Invalid id"
// @Failure 401   "Unauthorized"
// @Failure 403   "Permission denied"
// @Failure 500   "Other err"
// @Router /forms/{formid} [patch]
func (h *formHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(formUpdRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := &models.Form{Id: c.Param("formid")}

		if request.Title != nil {
			modelBL.Title = *request.Title
		}

		if request.Description != nil {
			modelBL.Description = *request.Description
		}

		updatedform, err := h.formUC.Update(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, formBLToResponse(updatedform))
	}
}

// GetById godoc
// @Summary Get form
// @Description Get form by id
// @Tags Forms
// @Security JWTToken
// @Param formid path string true "form id"
// @Success 200 {object} formResponse "Found"
// @Failure 204 {object} formResponse "No such form"
// @Failure 400   "Invalid id"
// @Failure 401   "Unauthorized"
// @Failure 500   "Other err"
// @Router /forms/{formid} [get]
func (h *formHandlers) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		foundform, err := h.formUC.GetById(c, c.Param("formid"))
		if err != nil {
			if err == errs.ErrContentNotFound {
				c.JSON(http.StatusNoContent, &formResponse{})
				return
			}

			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, formBLToResponse(foundform))
	}
}

func formCreatRequestToBL(dto *formCreatRequest) *models.Form {
	return &models.Form{
		Title:       dto.Title,
		Description: dto.Description,
	}
}

func formBLToResponse(modelBL *models.Form) *formResponse {
	return &formResponse{
		Id:          modelBL.Id,
		User_id:     modelBL.User_id,
		Title:       modelBL.Title,
		Description: modelBL.Description,
	}
}

func formsBLToResponse(forms []*models.Form) []*formResponse {
	if forms == nil {
		return nil
	}

	res := make([]*formResponse, len(forms))

	for i, f := range forms {
		res[i] = formBLToResponse(f)
	}

	return res
}
