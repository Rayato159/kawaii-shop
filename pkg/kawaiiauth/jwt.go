package kawaiiauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

type IKawaiiAuth interface {
	SignToken() string
}

type kawaiiAuth struct {
	mapClaims *kawaiiMapClaims
	cfg       config.IJwtConfig
}

type kawaiiMapClaims struct {
	Claims any `json:"claims"`
	jwt.RegisteredClaims
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *kawaiiAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*kawaiiMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &kawaiiMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check sign algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token error: %v", err)
		}
	}

	// Check type and return
	if claims, ok := token.Claims.(*kawaiiMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims any, exp int64) string {
	obj := &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewKawaiiAuth(tokenType TokenType, cfg config.IJwtConfig, claims any) (IKawaiiAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims any) IKawaiiAuth {
	return &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessTokenExpires()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims any) IKawaiiAuth {
	return &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshTokenExpires()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
