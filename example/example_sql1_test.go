package example

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL
	"github.com/stretchr/testify/assert"
	"sql-mapper/endpoint"
	"sql-mapper/reader"
	"sql-mapper/store"
	"testing"
)

func TestReadMapper(t *testing.T) {
	t.Run("mapper read test", func(t *testing.T) {
		body, e := reader.ReadMapperFile("./mapper/sql1.xml")
		if e != nil {
			fmt.Println(e)
		}

		fmt.Println(body)
	})
}

func TestRegister(t *testing.T) {
	t.Run("register test", func(t *testing.T) {
		register, err := store.Register("basic query", "./mapper/sql1.xml")
		assert.Nil(t, err)
		assert.NotNil(t, register)
	})
}

func TestNewQueryClient(t *testing.T) {
	db := sqlx.MustConnect(dbDriver, datasourceName)
	defer db.Close()

	var accounts []accountDb
	client, _ := endpoint.NewQueryClient(db, "identifier", "./mapper/sql1.xml")
	_ = client.Get(context.Background(), "allUsers", &accounts, map[string]any{})
	fmt.Println(accounts)
}
