package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (ac *apiConfig) createJwt(userId int, expiresIn time.Duration) (string, error) {
	issuedAt := jwt.NewNumericDate(time.Now().UTC())
	expiredAt := jwt.NewNumericDate(time.Now().UTC().Add(expiresIn))
	issuer := "chirpy"
	subject := fmt.Sprintf("%d", userId)

	fmt.Printf("%v\n", expiredAt)

	tok := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  issuedAt,
			ExpiresAt: expiredAt,
			Subject:   subject,
		})

	jwt, err := tok.SignedString([]byte(ac.jwtSecret))
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (ac *apiConfig) validateJwt(tok string) (int, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tok, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(ac.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	subject, err := token.Claims.(*jwt.RegisteredClaims).GetSubject()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
