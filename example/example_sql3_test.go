package example

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sql-mapper/endpoint"
	"sql-mapper/reader"
	"testing"
)

func TestReadMapperFile(t *testing.T) {
	raw, err := reader.ReadMapperFile("./mapper/sql3.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(raw)
}

func TestNewQueryClient3(t *testing.T) {
	_, err := endpoint.NewQueryClient(db, "identifier-TestNewQueryClient", "./mapper/sql3.xml")
	assert.Nil(t, err)
}
