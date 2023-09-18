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

func TestMain(m *testing.M) {
	// before the test
	db = sqlx.MustConnect(dbDriver, datasourceName)
	c1, _ := endpoint.NewQueryClient(db, "identifier1", "./mapper/sql1.xml")
	client1 = c1
	c2, _ := endpoint.NewQueryClient(db, "identifier2", "./mapper/sql2.xml")
	client2 = c2

	// run tests
	exitCode := m.Run()

	// after the test
	_ = db.Close()
	os.Exit(exitCode)
}
