package example

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
