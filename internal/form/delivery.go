package form

import "github.com/gin-gonic/gin"

// Form HTTP Handlers interface
type Handlers interface {
	Create() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Update() gin.HandlerFunc
	GetById() gin.HandlerFunc
	GetByUser() gin.HandlerFunc
}
