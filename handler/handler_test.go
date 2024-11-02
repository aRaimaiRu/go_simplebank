package handler_test

import (
	"fmt"
	db "go_simplebank/db/sqlc"
	"go_simplebank/handler"
	"go_simplebank/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *handler.Server {
	config := util.Config{
		DBDriver:            "postgres",
		DBSource:            "postgresql",
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := handler.NewServer(config, store)
	require.NoError(t, err)
	return server
}
func TestMain(m *testing.M) {
	fmt.Println("TestMain api")
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
