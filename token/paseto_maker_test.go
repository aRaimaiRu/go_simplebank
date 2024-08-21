package token_test

import (
	"go_simplebank/token"
	"go_simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	IssuedAt := time.Now()
	ExpiredAt := IssuedAt.Add(duration)

	token, _, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, ExpiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredpasetoToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	createdToken, _, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, createdToken)

	payload, err := maker.VerifyToken(createdToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidpasetoToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("invalid-token")
	require.Error(t, err)
	require.EqualError(t, err, token.ErrTokenInvalid.Error())
	require.Nil(t, payload)
}
