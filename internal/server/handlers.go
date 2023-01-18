package server

import (
	"net/http"
	arepo "quiz-app/internal/answer/repo"
	auc "quiz-app/internal/answer/usecase"
	authh "quiz-app/internal/auth/delivery/http"
	authrepo "quiz-app/internal/auth/repo"
	authuc "quiz-app/internal/auth/usecase"
	fh "quiz-app/internal/form/delivery/http"
	frepo "quiz-app/internal/form/repo"
	fuc "quiz-app/internal/form/usecase"
	pah "quiz-app/internal/poolanswer/delivery/http"
	parepo "quiz-app/internal/poolanswer/repo"
	pauc "quiz-app/internal/poolanswer/usecase"
	qh "quiz-app/internal/question/delivery/http"
	qrepo "quiz-app/internal/question/repo"
	quc "quiz-app/internal/question/usecase"
	jwtgo "quiz-app/pkg/jwter/impl"

	_ "quiz-app/docs"

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
