package example

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"sql-mapper/endpoint"
	"testing"
)

func TestFullDMLNonTx(t *testing.T) {
	// db connection
	db := sqlx.MustConnect(dbDriver, datasourceName)
	defer db.Close()

	// create context
	ctx := context.Background()

	// get QueryClient
	client, _ := endpoint.NewQueryClient(db, "identifier", "./mapper/sql2.xml")

	t.Run("full-test-non-tx", func(t *testing.T) {
		_, err := client.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(accounts))

		err = client.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		err = client.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(accounts))

		account := new(accountDb)
		err = client.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "test1", account.Name)
		assert.Equal(t, "test1-ID", account.UserId)
		assert.Equal(t, "test1-PW", account.Password)
		assert.True(t, account.Id > 0)

		nums, err := client.Update(ctx, "updateUserNameForOneUser", map[string]any{
			"user_name": "update-name",
			"user_id":   "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), nums)

		err = client.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "update-name", account.Name)
		nums, err = client.Delete(ctx, "deleteOneUser", map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), nums)
	})
}

func TestFullDmlWithTx(t *testing.T) {
	// db connection
	db := sqlx.MustConnect(dbDriver, datasourceName)
	defer db.Close()

	// create context
	ctx := context.Background()

	// get QueryClient
	client, _ := endpoint.NewQueryClient(db, "identifier", "./mapper/sql2.xml")

	t.Run("insert-rollback-test", func(t *testing.T) {
		// reset
		_, err := client.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		tx, err := client.BeginTx(ctx)
		assert.Nil(t, err)
		err = client.InsertOneTx(ctx, tx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)
		err = client.RollbackTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(accounts))
	})

	t.Run("insert-commit-test", func(t *testing.T) {
		// reset
		_, err := client.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		tx, err := client.BeginTx(ctx)
		assert.Nil(t, err)
		err = client.InsertOneTx(ctx, tx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)
		err = client.CommitTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(accounts))
	})

	t.Run("delete-rollback-test", func(t *testing.T) {
		// reset
		_, err := client.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client.BeginTx(ctx)
		assert.Nil(t, err)
		nums, err := client.DeleteTx(ctx, tx, "deleteOneUser", map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), nums)
		err = client.RollbackTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(accounts))
	})

	t.Run("delete-commit-test", func(t *testing.T) {
		// reset
		_, err := client.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client.BeginTx(ctx)
		assert.Nil(t, err)
		nums, err := client.DeleteTx(ctx, tx, "deleteOneUser", map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), nums)
		err = client.CommitTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(accounts))
	})

	t.Run("update-rollback-test", func(t *testing.T) {
		// reset
		_, err := client.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client.BeginTx(ctx)
		assert.Nil(t, err)
		rowNums, err := client.UpdateTx(ctx, tx, "updateUserNameForOneUser", map[string]any{
			"user_name": "update-name",
			"user_id":   "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), rowNums)
		err = client.RollbackTx(ctx, tx)
		assert.Nil(t, err)

		account := new(accountDb)
		err = client.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "test1", account.Name)
	})

	t.Run("update-commit-test", func(t *testing.T) {
		// reset
		_, err := client.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client.BeginTx(ctx)
		assert.Nil(t, err)
		rowNums, err := client.UpdateTx(ctx, tx, "updateUserNameForOneUser", map[string]any{
			"user_name": "update-name",
			"user_id":   "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), rowNums)
		err = client.CommitTx(ctx, tx)
		assert.Nil(t, err)

		account := new(accountDb)
		err = client.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "update-name", account.Name)
	})
}
