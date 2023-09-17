package example

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFullDMLNonTx(t *testing.T) {
	// create context
	ctx := context.Background()

	_, err := client2.Delete(ctx, "fullDelete", map[string]any{})
	assert.Nil(t, err)

	var accounts []accountDb
	err = client2.Get(ctx, "allUsers", &accounts, map[string]any{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(accounts))

	err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
		"user_name": "test1",
		"user_id":   "test1-ID",
		"password":  "test1-PW",
	})
	assert.Nil(t, err)

	err = client2.Get(ctx, "allUsers", &accounts, map[string]any{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(accounts))

	account := new(accountDb)
	err = client2.GetOne(ctx, "specificUser", account, map[string]any{
		"user_id": "test1-ID",
	})
	assert.Nil(t, err)
	assert.Equal(t, "test1", account.Name)
	assert.Equal(t, "test1-ID", account.UserId)
	assert.Equal(t, "test1-PW", account.Password)
	assert.True(t, account.Id > 0)

	nums, err := client2.Update(ctx, "updateUserNameForOneUser", map[string]any{
		"user_name": "update-name",
		"user_id":   "test1-ID",
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), nums)

	err = client2.GetOne(ctx, "specificUser", account, map[string]any{
		"user_id": "test1-ID",
	})
	assert.Nil(t, err)
	assert.Equal(t, "update-name", account.Name)
	nums, err = client2.Delete(ctx, "deleteOneUser", map[string]any{
		"user_id": "test1-ID",
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), nums)
}

func TestFullDmlWithTx(t *testing.T) {
	// create context
	ctx := context.Background()

	t.Run("insert-rollback-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		tx, err := client2.BeginTx(ctx)
		assert.Nil(t, err)
		err = client2.InsertOneTx(ctx, tx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)
		err = client2.RollbackTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client2.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(accounts))
	})

	t.Run("insert-commit-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		tx, err := client2.BeginTx(ctx)
		assert.Nil(t, err)
		err = client2.InsertOneTx(ctx, tx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)
		err = client2.CommitTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client2.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(accounts))
	})

	t.Run("delete-rollback-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client2.BeginTx(ctx)
		assert.Nil(t, err)
		nums, err := client2.DeleteTx(ctx, tx, "deleteOneUser", map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), nums)
		err = client2.RollbackTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client2.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(accounts))
	})

	t.Run("delete-commit-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client2.BeginTx(ctx)
		assert.Nil(t, err)
		nums, err := client2.DeleteTx(ctx, tx, "deleteOneUser", map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), nums)
		err = client2.CommitTx(ctx, tx)
		assert.Nil(t, err)

		var accounts []accountDb
		err = client2.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(accounts))
	})

	t.Run("update-rollback-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client2.BeginTx(ctx)
		assert.Nil(t, err)
		rowNums, err := client2.UpdateTx(ctx, tx, "updateUserNameForOneUser", map[string]any{
			"user_name": "update-name",
			"user_id":   "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), rowNums)
		err = client2.RollbackTx(ctx, tx)
		assert.Nil(t, err)

		account := new(accountDb)
		err = client2.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "test1", account.Name)
	})

	t.Run("update-commit-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		// input sample
		err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
			"user_name": "test1",
			"user_id":   "test1-ID",
			"password":  "test1-PW",
		})
		assert.Nil(t, err)

		tx, err := client2.BeginTx(ctx)
		assert.Nil(t, err)
		rowNums, err := client2.UpdateTx(ctx, tx, "updateUserNameForOneUser", map[string]any{
			"user_name": "update-name",
			"user_id":   "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), rowNums)
		err = client2.CommitTx(ctx, tx)
		assert.Nil(t, err)

		account := new(accountDb)
		err = client2.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "update-name", account.Name)
	})
}
