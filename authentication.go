package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokens struct {
	accessToken  string
	refreshToken string
}

func (ac *apiConfig) createJwt(userId int) (*JwtTokens, error) {
	issuedAt := jwt.NewNumericDate(time.Now().UTC())

	expireRefresh := time.Duration(time.Hour * 24 * 60)

	expiredAtAccess := jwt.NewNumericDate(time.Now().UTC().Add(time.Hour))
	expiredAtRefresh := jwt.NewNumericDate(time.Now().UTC().Add(expireRefresh))
	issuerAccess := "chirpy-access"
	issuerRefresh := "chirpy-refresh"
	subject := fmt.Sprintf("%d", userId)

	accessTok := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuerAccess,
			IssuedAt:  issuedAt,
			ExpiresAt: expiredAtAccess,
			Subject:   subject,
		})

	refreshTok := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuerRefresh,
			IssuedAt:  issuedAt,
			ExpiresAt: expiredAtRefresh,
			Subject:   subject,
		})

	jwtAccess, err := accessTok.SignedString([]byte(ac.jwtSecret))
	if err != nil {
		return nil, err
	}

	jwtRefresh, err := refreshTok.SignedString([]byte(ac.jwtSecret))
	if err != nil {
		return nil, err
	}

	t := JwtTokens{accessToken: jwtAccess, refreshToken: jwtRefresh}

	return &t, nil
}

// i have no idea has this should work??
// why do i have to provide claims? are there filled after or not???
func (ac *apiConfig) validateJwt(tok string) (int, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tok, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(ac.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	_claims := token.Claims.(*jwt.RegisteredClaims)

	if _claims.Issuer != "chirpy-access" {
		return 0, errors.New("wrong issuer")
	}

	subject, err := _claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (ac *apiConfig) validateRefreshJwt(tok string) (int, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tok, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(ac.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	_claims := token.Claims.(*jwt.RegisteredClaims)

	if _claims.Issuer != "chirpy-refresh" {
		return 0, errors.New("wrong issuer")
	}

	subject, err := _claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
