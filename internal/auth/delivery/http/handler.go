package http

import (
	"net/http"
	"quizapp/internal/auth"
	"quizapp/models"
	"quizapp/pkg/errs"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type GetResponse struct {
	Login string `json:"login"`
}

type authHandlers struct {
	authUC auth.UseCase
}

func NewAuthHandlers(authUC auth.UseCase) auth.Handlers {
	return &authHandlers{
		authUC: authUC,
	}
}

// SignUp godoc
// @Summary Sign up user
// @Description Sign up new user in the system for further sign in and use
// @Tags Auth
// @Accept json
// @Param credentials body AuthRequest true "user credentials"
// @Success 201 {object} SignUpResponse
// @Failure 400   "Invalid json"
// @Failure 403   "Permission denied"
// @Failure 409   "User with such login exists"
// @Failure 500   "Other err"
// @Router /auth/signup [post]
func (h *authHandlers) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(AuthRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		signedupuser, err := h.authUC.SignUp(c.Request.Context(), requestToBL(request))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusCreated, blToSignUpResponse(signedupuser))
	}
}

// SignIn godoc
// @Summary Sign in user
// @Description Sign in user by login and password to the system for further use
// @Tags Auth
// @Accept json
// @Param credentials body AuthRequest true "user credentials"
// @Success 201 {object} SignInResponse
// @Failure 400   "Invalid json"
// @Failure 401   "Invalid login or password"
// @Failure 403   "Permission denied"
// @Failure 500   "Other err"
// @Router /auth/signin [post]
func (h *authHandlers) SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(AuthRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := h.authUC.SignIn(c.Request.Context(), requestToBL(request))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, SignInResponse{Token: *token})
	}
}

// GetById godoc
// @Summary Get user login
// @Description Get user login by id
// @Tags Auth
// @Security JWTToken
// @Param id path string true "user id"
// @Success 201 {object} GetResponse
// @Success 204   "No such user"
// @Failure 400   "Invalid params"
// @Failure 500   "Other err"
// @Router /users/{id} [get]
func (h *authHandlers) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := h.authUC.GetById(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, GetResponse{user.Login})
	}
}

func requestToBL(request *AuthRequest) *models.User {
	return &models.User{
		Login:    request.Login,
		Password: request.Password,
	}
}

func blToSignUpResponse(user *models.User) *SignUpResponse {
	return &SignUpResponse{
		Id:       user.Id,
		Login:    user.Login,
		Password: user.Password,
	}
}
