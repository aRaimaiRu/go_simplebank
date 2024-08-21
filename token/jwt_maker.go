package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Payload
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(username, duration)
	claims := JwtClaims{
		Payload: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrTokenInvalid
	}
	payload, ok := jwtToken.Claims.(*JwtClaims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return &payload.Payload, nil
}

func EncodeUsernameToInt64(username string) (int64, error) {
	usernameBytes := []byte(username)
	userID := int64(usernameBytes[0])
	return userID, nil
}

// func (p *jwtClaims) GetAudience() (jwt.ClaimStrings, error) {
// 	return jwt.ClaimStrings{p.Username}, nil
// }

// func (p *jwtClaims) GetExpirationTime() (*jwt.NumericDate, error) {
// 	return &jwt.NumericDate{p.ExpiredAt}, nil
// }
