package example

import (
	"context"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL
	"testing"
)

func init() {

}

func TestNewQueryClient(t *testing.T) {
	var accounts []accountDb
	_ = client1.Get(context.Background(), "allUsers", &accounts, map[string]any{})
	fmt.Println(accounts)
}
