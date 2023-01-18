package poolanswer

import "github.com/gin-gonic/gin"

// Question HTTP Handlers interface
type Handlers interface {
	Create() gin.HandlerFunc
	GetByFormId() gin.HandlerFunc
	GetByPoolAnswerId() gin.HandlerFunc
}
