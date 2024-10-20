package handler_test

import (
	"fmt"
	"go_simplebank/api"
	db "go_simplebank/db/sqlc"
	"go_simplebank/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *api.Server {
	config := util.Config{
		DBDriver:            "postgres",
		DBSource:            "postgresql",
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := api.NewServer(config, store)
	require.NoError(t, err)
	return server
}
func TestMain(m *testing.M) {
	fmt.Println("TestMain api")
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
