package question

import "github.com/gin-gonic/gin"

// Question HTTP Handlers interface
type Handlers interface {
	Create() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Update() gin.HandlerFunc
	GetByFormId() gin.HandlerFunc
}
