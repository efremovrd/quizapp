package http

import (
	"quiz-app/internal/question"

	"github.com/gin-gonic/gin"
)

// Map question routes
func MapQuestionRoutes(questionGroup *gin.RouterGroup, h question.Handlers) {
	questionGroup.POST("", h.Create())
	questionGroup.GET("", h.GetByFormId())
	questionGroup.PUT("/:questionid", h.Update())
	questionGroup.DELETE("/:questionid", h.Delete())
}
