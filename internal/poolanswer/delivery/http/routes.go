package http

import (
	"quiz-app/internal/poolanswer"

	"github.com/gin-gonic/gin"
)

// Map answers routes
func MapPARoutes(answersGroup *gin.RouterGroup, h poolanswer.Handlers) {
	answersGroup.POST("", h.Create())
	answersGroup.GET("", h.GetByFormId())
	answersGroup.GET("/:poolanswerid", h.GetByPoolAnswerId())
}
