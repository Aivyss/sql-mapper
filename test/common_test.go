package test

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL
	"os"
	"sql-mapper/context"
	"sql-mapper/test/helper"
	"testing"
)

const dbDriver = "postgres"
const datasourceName = "user=test1 password=test1 dbname=test1 sslmode=disable"

func TestMain(m *testing.M) {
	// before the test
	db := sqlx.MustConnect(dbDriver, datasourceName)
	_, err := context.Bootstrap(db).InitByXml("./setting/settings.xml")
	helper.DoPanicIfNotNil(err)

	// run tests
	exitCode := m.Run()

	// after the test
	_ = db.Close()
	os.Exit(exitCode)
}
