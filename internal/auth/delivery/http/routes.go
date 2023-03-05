package http

import (
	"quizapp/internal/auth"

	"github.com/gin-gonic/gin"
)

// Map auth routes
func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers) {
	authGroup.POST("/signup", h.SignUp())
	authGroup.POST("/signin", h.SignIn())
}
