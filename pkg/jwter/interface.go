package jwter

import "quiz-app/models"

type JWTer interface {
	// Returns generated token & nil, if generated.
	// Returns nil & some err else.
	GenerateJWTToken(user *models.User) (*string, error)

	// Returns user model & nil, if parsed.
	// Returns nil & some err else.
	ParseToken(access_token string) (*models.User, error)
}
