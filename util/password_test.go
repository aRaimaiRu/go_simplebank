package util_test

import (
	"go_simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	checkPassword := util.CheckPasswordHash(password, hashedPassword)
	require.NoError(t, checkPassword)

	wrongPassword := util.RandomString(6)
	checkPassword = util.CheckPasswordHash(wrongPassword, hashedPassword)
	require.NoError(t, checkPassword)
}
