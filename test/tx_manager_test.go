package test

import (
	"context"
	"errors"
	lctx "github.com/aivyss/sql-mapper/context"
	"github.com/aivyss/sql-mapper/test/helper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactionManager1(t *testing.T) {
	// create context
	ctx := context.Background()
	txManager := lctx.GetApplicationContext().GetTxManager()

	client2, err := lctx.GetApplicationContext().GetQueryClient("identifier2")
	helper.DoPanicIfNotNil(err)

	t.Run("insert-rollback-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		_ = txManager.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			assert.Nil(t, err)
			err = client2.InsertOneTx(ctx, tx, "saveOneUser", map[string]any{
				"user_name": "test1",
				"user_id":   "test1-ID",
				"password":  "test1-PW",
			})
			assert.Nil(t, err)

			return errors.New("rollback test err")
		})

		var accounts []accountDb
		err = client2.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(accounts))
	})

	t.Run("insert-commit-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		_ = txManager.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			assert.Nil(t, err)
			err = client2.InsertOneTx(ctx, tx, "saveOneUser", map[string]any{
				"user_name": "test1",
				"user_id":   "test1-ID",
				"password":  "test1-PW",
			})
			assert.Nil(t, err)

			return nil
		})

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

		_ = txManager.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			nums, err := client2.DeleteTx(ctx, tx, "deleteOneUser", map[string]any{
				"user_id": "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), nums)

			return errors.New("rollback test err")
		})

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

		_ = txManager.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			nums, err := client2.DeleteTx(ctx, tx, "deleteOneUser", map[string]any{
				"user_id": "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), nums)

			return nil
		})

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

		_ = txManager.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			rowNums, err := client2.UpdateTx(ctx, tx, "updateUserNameForOneUser", map[string]any{
				"user_name": "update-name",
				"user_id":   "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), rowNums)

			return errors.New("rollback test err")
		})

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

		_ = txManager.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			rowNums, err := client2.UpdateTx(ctx, tx, "updateUserNameForOneUser", map[string]any{
				"user_name": "update-name",
				"user_id":   "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), rowNums)

			return nil
		})

		account := new(accountDb)
		err = client2.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "update-name", account.Name)
	})
}

func TestTransactionManager2(t *testing.T) {
	// create context
	ctx := context.Background()
	txManager := lctx.GetApplicationContext().GetTxManager()

	client2, err := lctx.GetApplicationContext().GetQueryClient("identifier2")
	helper.DoPanicIfNotNil(err)

	t.Run("insert-rollback-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		_ = txManager.Txx(ctx, func(ctx context.Context) error {
			assert.Nil(t, err)
			err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
				"user_name": "test1",
				"user_id":   "test1-ID",
				"password":  "test1-PW",
			})
			assert.Nil(t, err)

			return errors.New("rollback test err")
		})

		var accounts []accountDb
		err = client2.Get(ctx, "allUsers", &accounts, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(accounts))
	})

	t.Run("insert-commit-test", func(t *testing.T) {
		// reset
		_, err := client2.Delete(ctx, "fullDelete", nil)
		assert.Nil(t, err)

		_ = txManager.Txx(ctx, func(ctx context.Context) error {
			assert.Nil(t, err)
			err = client2.InsertOne(ctx, "saveOneUser", map[string]any{
				"user_name": "test1",
				"user_id":   "test1-ID",
				"password":  "test1-PW",
			})
			assert.Nil(t, err)

			return nil
		})

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

		_ = txManager.Txx(ctx, func(ctx context.Context) error {
			nums, err := client2.Delete(ctx, "deleteOneUser", map[string]any{
				"user_id": "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), nums)

			return errors.New("rollback test err")
		})

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

		_ = txManager.Txx(ctx, func(ctx context.Context) error {
			nums, err := client2.Delete(ctx, "deleteOneUser", map[string]any{
				"user_id": "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), nums)

			return nil
		})

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

		_ = txManager.Txx(ctx, func(ctx context.Context) error {
			rowNums, err := client2.Update(ctx, "updateUserNameForOneUser", map[string]any{
				"user_name": "update-name",
				"user_id":   "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), rowNums)

			return errors.New("rollback test err")
		})

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

		_ = txManager.Txx(ctx, func(ctx context.Context) error {
			rowNums, err := client2.Update(ctx, "updateUserNameForOneUser", map[string]any{
				"user_name": "update-name",
				"user_id":   "test1-ID",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(1), rowNums)

			return nil
		})

		account := new(accountDb)
		err = client2.GetOne(ctx, "specificUser", account, map[string]any{
			"user_id": "test1-ID",
		})
		assert.Nil(t, err)
		assert.Equal(t, "update-name", account.Name)
	})
}
