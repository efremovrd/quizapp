package http

import (
	"quizapp/internal/form"

	"github.com/gin-gonic/gin"
)

// Map form routes
func MapFormRoutes(formGroup *gin.RouterGroup, h form.Handlers) {
	formGroup.POST("", h.Create())
	formGroup.DELETE("/:formid", h.Delete())
	formGroup.PATCH("/:formid", h.Update())
	formGroup.GET("", h.GetByUser())
	formGroup.GET("/:formid", h.GetById())
}
