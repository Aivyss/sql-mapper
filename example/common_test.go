package example

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL
	"os"
	"sql-mapper/endpoint"
	"testing"
)

const dbDriver = "postgres"
const datasourceName = "user=test1 password=test1 dbname=test1 sslmode=disable"

var db *sqlx.DB
var client1 endpoint.QueryClient
var client2 endpoint.QueryClient
var client3 endpoint.QueryClient

func TestMain(m *testing.M) {
	// before the test
	db = sqlx.MustConnect(dbDriver, datasourceName)
	client1, _ = endpoint.NewQueryClient(db, "identifier1", "./mapper/sql1.xml")
	client2, _ = endpoint.NewQueryClient(db, "identifier2", "./mapper/sql2.xml")
	client3, _ = endpoint.NewQueryClient(db, "identifier3", "./mapper/sql3.xml")

	// run tests
	exitCode := m.Run()

	// after the test
	_ = db.Close()
	os.Exit(exitCode)
}
