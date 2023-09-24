package test

import (
	"context"
	"fmt"
	ctx "sql-mapper/context"
	"sql-mapper/test/helper"
	"testing"
)

func init() {

}

func TestNewQueryClient(t *testing.T) {
	client1, err := ctx.GetApplicationContext().GetQueryClient("identifier1")
	helper.DoPanicIfNotNil(err)

	var accounts []accountDb
	_ = client1.Get(context.Background(), "allUsers", &accounts, map[string]any{})
	fmt.Println(accounts)
}
