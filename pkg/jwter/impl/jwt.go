package jwtgo

import (
	"quiz-app/models"
	"quiz-app/pkg/jwter"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	User *models.User
	jwt.StandardClaims
}

type jwtgo struct {
	key string
}

func NewJWTGO(secret_key string) jwter.JWTer {
	return &jwtgo{secret_key}
}

func (j *jwtgo) GenerateJWTToken(user *models.User) (*string, error) {
	claims := &claims{
		User:           user,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token_string, err := token.SignedString([]byte(j.key))
	if err != nil {
		return nil, err
	}

	return &token_string, nil
}

func (j *jwtgo) ParseToken(access_token string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(access_token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.key), nil
	})

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, err
}
