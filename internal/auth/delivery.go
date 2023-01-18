package auth

import "github.com/gin-gonic/gin"

// Auth HTTP Handlers interface
type Handlers interface {
	SignUp() gin.HandlerFunc
	SignIn() gin.HandlerFunc
	GetById() gin.HandlerFunc
}
