package utils

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/golang-jwt/jwt/v5"
)

type IKawaiiAuth interface {
	SignToken()
	ParseToken()
}

type KawaiiAuth struct {
}

type TokenMapClaims struct {
	Claims any
	jwt.RegisteredClaims
}

func NewKawaiiAuth(cfg config.IJwtConfig, claims any) {

}
