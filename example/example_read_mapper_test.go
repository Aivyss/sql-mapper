package example

import (
	"fmt"
	"sql-mapper/reader"
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
