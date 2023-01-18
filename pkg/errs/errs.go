package errs

import (
	"errors"
	"net/http"
)

var (
	ErrContentNotFound    = errors.New("content not found")
	ErrUnauthorized       = errors.New("user unathorized")
	ErrForbidden          = errors.New("user not an owner")
	ErrInvalidContent     = errors.New("invalid content")
	ErrLoginExists        = errors.New("login already exists")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrInvalidPassword    = errors.New("invalid password")
)

func MatchHttpErr(err error) int {
	if err == ErrContentNotFound {
		return http.StatusNoContent
	}

	if err == ErrInvalidContent {
		return http.StatusBadRequest
	}

	if err == ErrForbidden {
		return http.StatusForbidden
	}

	if err == ErrUnauthorized ||
		err == ErrInvalidAccessToken ||
		err == ErrInvalidPassword {
		return http.StatusUnauthorized
	}

	if err == ErrLoginExists {
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}
