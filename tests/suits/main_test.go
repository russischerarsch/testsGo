package suits

import (
	"context"
	"log"
	"os"
	"testing"
	"tgtest/tests"
)

var Brokers []string
var Ctx = context.Background()
var Dsn string

func TestMain(m *testing.M) {
	teardown, brokers, err := tests.SetupContainers()
	if err != nil {
		log.Fatalf("failed to set up containers")
	}
	defer teardown()
	connStr := os.Getenv("POSTGRES_CONN")
	Dsn = connStr
	Brokers = brokers
	os.Exit(m.Run())
}
