package server

import (
	"net/http"
	arepo "quizapp/internal/answer/repo"
	auc "quizapp/internal/answer/usecase"
	authh "quizapp/internal/auth/delivery/http"
	authrepo "quizapp/internal/auth/repo"
	authuc "quizapp/internal/auth/usecase"
	fh "quizapp/internal/form/delivery/http"
	frepo "quizapp/internal/form/repo"
	fuc "quizapp/internal/form/usecase"
	pah "quizapp/internal/poolanswer/delivery/http"
	parepo "quizapp/internal/poolanswer/repo"
	pauc "quizapp/internal/poolanswer/usecase"
	qh "quizapp/internal/question/delivery/http"
	qrepo "quizapp/internal/question/repo"
	quc "quizapp/internal/question/usecase"
	jwtgo "quizapp/pkg/jwter/impl"

	_ "quizapp/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) MapHandlers() error {
	jwter := jwtgo.NewJWTGO(s.cfg.Server.JwtSecretKey)

	aRepo := arepo.NewAnswerRepo(s.db)
	authRepo := authrepo.NewAuthRepo(s.db)
	fRepo := frepo.NewFormRepo(s.cfg.Server.CtxUserKey, s.db)
	qRepo := qrepo.NewQuestionRepo(s.db)
	paRepo := parepo.NewPoolAnswerRepo(s.db)

	paUC := pauc.NewPoolAnswerUseCase(paRepo, aRepo, fRepo)
	aUC := auc.NewAnswerUseCase(aRepo, fRepo, paRepo)
	qUC := quc.NewQuestionUseCase(qRepo, fRepo)
	authUC := authuc.NewAuthUseCase(authRepo, jwter)
	fUC := fuc.NewFormUseCase(fRepo, s.cfg.Server.CtxUserKey)

	authH := authh.NewAuthHandlers(authUC)
	middleware := authh.NewAuthMiddleware(authUC, s.cfg.Server.CtxUserKey)
	fH := fh.NewFormHandlers(fUC, s.cfg.Server.CtxUserKey)
	qH := qh.NewQuestionHandlers(qUC, s.cfg.Server.CtxUserKey)
	aH := pah.NewAnswersHandlers(paUC, aUC, fUC, s.cfg.Server.CtxUserKey)

	s.router.GET("api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router.GET("api/v1/", func(c *gin.Context) { c.Redirect(http.StatusSeeOther, "/api/v1/docs/index.html") })

	auth := s.router.Group("api/v1/auth")
	authh.MapAuthRoutes(auth, authH)

	api := s.router.Group("/api", middleware)

	v1 := api.Group("/v1")

	v1.GET("/users/:id", authH.GetById())

	forms := v1.Group("/forms")
	fh.MapFormRoutes(forms, fH)

	questions := forms.Group("/:formid/questions")
	qh.MapQuestionRoutes(questions, qH)

	answers := forms.Group("/:formid/poolsanswer")
	pah.MapPARoutes(answers, aH)

	return nil
}
