package token_test

import (
	"go_simplebank/token"
	"go_simplebank/util"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	IssuedAt := time.Now()
	ExpiredAt := IssuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, ExpiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	createdToken, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, createdToken)

	payload, err := maker.VerifyToken(createdToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload := token.NewPayload(util.RandomOwner(), time.Minute)
	claims := token.JwtClaims{
		Payload: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	newToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(newToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrTokenInvalid.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("invalid-token")
	require.Error(t, err)
	require.EqualError(t, err, token.ErrTokenInvalid.Error())
	require.Nil(t, payload)
}
