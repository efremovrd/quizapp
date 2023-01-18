package http

import (
	"net/http"
	"quiz-app/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authMiddleware struct {
	authUC     auth.UseCase
	ctxUserKey string
}

func NewAuthMiddleware(authUC auth.UseCase, ctxUserKey string) gin.HandlerFunc {
	return (&authMiddleware{
		authUC:     authUC,
		ctxUserKey: ctxUserKey,
	}).Handle
}

func (m *authMiddleware) Handle(c *gin.Context) {
	authheader := c.GetHeader("Authorization")
	if authheader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerparts := strings.Split(authheader, " ")
	if len(headerparts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerparts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.authUC.ParseToken(c.Request.Context(), headerparts[1])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	founduser, err := m.authUC.GetById(c, user.Id)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if founduser.Login == user.Login && bcrypt.CompareHashAndPassword([]byte(founduser.Password), []byte(user.Password)) == nil {
		c.Set(m.ctxUserKey, user)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
