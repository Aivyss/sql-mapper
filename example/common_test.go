package example

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL
	"os"
	"sql-mapper/context"
	"testing"
)

const dbDriver = "postgres"
const datasourceName = "user=test1 password=test1 dbname=test1 sslmode=disable"

var db *sqlx.DB
var client1 context.QueryClient
var client2 context.QueryClient
var client3 context.QueryClient

func TestMain(m *testing.M) {
	// before the test
	db = sqlx.MustConnect(dbDriver, datasourceName)
	client1, _ = context.NewQueryClient("identifier1", "./mapper/sql1.xml", false)
	client2, _ = context.NewQueryClient("identifier2", "./mapper/sql2.xml", false)
	client3, _ = context.NewQueryClient("identifier3", "./mapper/sql3.xml", false)

	// run tests
	exitCode := m.Run()

	// after the test
	_ = db.Close()
	os.Exit(exitCode)
}
