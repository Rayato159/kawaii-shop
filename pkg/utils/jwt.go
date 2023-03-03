package utils

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/golang-jwt/jwt/v5"
)

type IKawaiiAuth interface {
	SignToken() *oauth.UserToken
	ParseToken() *kawaiiMapClaims
	RefreshToken() string
}

func (a *kawaiiAuth) SignToken() *oauth.UserToken {
	return a.token
}

func (a *kawaiiAuth) ParseToken() *kawaiiMapClaims {
	return a.mapClaims
}

func (a *kawaiiAuth) RefreshToken() string {
	return a.token.RefreshToken
}

type kawaiiAuth struct {
	token     *oauth.UserToken
	mapClaims *kawaiiMapClaims
}

type kawaiiMapClaims struct {
	claims any
	jwt.RegisteredClaims
}

func NewKawaiiAuth(cfg config.IJwtConfig, claims any) IKawaiiAuth {
	return &kawaiiAuth{
		mapClaims: &kawaiiMapClaims{
			claims: claims,
		},
	}
}
