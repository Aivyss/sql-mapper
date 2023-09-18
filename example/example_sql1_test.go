package example

import (
	"context"
	"fmt"
	"testing"
)

func init() {

}

func TestNewQueryClient(t *testing.T) {
	var accounts []accountDb
	_ = client1.Get(context.Background(), "allUsers", &accounts, map[string]any{})
	fmt.Println(accounts)
}
