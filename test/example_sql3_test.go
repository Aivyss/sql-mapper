package test

import (
	"context"
	"fmt"
	lctx "github.com/aivyss/sql-mapper/context"
	"github.com/aivyss/sql-mapper/entity"
	"github.com/aivyss/sql-mapper/reader/xml"
	"github.com/aivyss/sql-mapper/test/helper"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestReadMapperFile(t *testing.T) {
	raw, err := xml.ReadMapperFile("./mapper/sql3.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(raw)
}

func TestNewQueryClient3(t *testing.T) {
	_, err := lctx.NewQueryClient("identifier-TestNewQueryClient", "./mapper/sql3.xml")
	assert.Nil(t, err)
}

func TestQueries(t *testing.T) {
	ctx := context.Background()
	client2, err := lctx.GetApplicationContext().GetQueryClient("identifier2")
	helper.DoPanicIfNotNil(err)
	client3, err := lctx.GetApplicationContext().GetQueryClient("identifier3")
	helper.DoPanicIfNotNil(err)

	t.Run("get select", func(t *testing.T) {
		// reset data
		_, err := client2.Delete(ctx, "fullDelete", map[string]any{})
		assert.Nil(t, err)

		err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		// get test with predicates
		dest := new(accountDb)
		randInt1 := rand.Intn(10)
		randInt2 := rand.Intn(10)
		err = client3.GetOne(ctx, "specificUser", dest, map[string]any{
			"user_id": "test1-ID",
		}, func() []*entity.Condition {
			var result []*entity.Condition

			if randInt1 > 4 {

				result = append(result, &entity.Condition{
					PartName: "condition1",
					CaseName: "case2",
				})
			} else {
				result = append(result, &entity.Condition{
					PartName: "condition1",
					CaseName: "case1",
				})
			}

			if randInt2 > 4 {
				result = append(result, &entity.Condition{
					PartName: "condition2",
					CaseName: "case4",
				})
			} else {
				result = append(result, &entity.Condition{
					PartName: "condition2",
					CaseName: "case3",
				})

			}

			return result
		})
		assert.Nil(t, err)
		assert.NotEqual(t, accountDb{}, dest)
		name := dest.Name
		password := dest.Password

		if name == "" && password == "" {
			assert.True(t, randInt1 > 4)
			assert.True(t, randInt2 > 4)
		} else if name == "" {
			assert.True(t, randInt1 > 4)
			assert.False(t, randInt2 > 4)
		} else if password == "" {
			assert.False(t, randInt1 > 4)
			assert.True(t, randInt2 > 4)
		} else {
			assert.False(t, randInt1 > 4)
			assert.False(t, randInt2 > 4)
		}
	})
}
