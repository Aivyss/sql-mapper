package example

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sql-mapper/incubate/endpoint"
	"sql-mapper/incubate/reader"
	"testing"
)

func TestReadMapperFile(t *testing.T) {
	raw, err := reader.ReadMapperFile("./mapper/sql3.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(raw)
}

func TestNewQueryClient(t *testing.T) {
	_, err := endpoint.NewQueryClient(db, "identifier-TestNewQueryClient", "./mapper/sql3.xml")
	assert.Nil(t, err)
}
