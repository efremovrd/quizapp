package jwter

type JWTer interface {
	// Returns generated token & nil, if generated.
	// Returns nil & some err else.
	GenerateJWTToken(id, login string) (*string, error)

	// Returns user model & nil, if parsed.
	// Returns nil & some err else.
	ParseToken(access_token string) (*string, *string, error)
}
