package endpoint

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"sql-mapper/enum"
	"sql-mapper/errors"
)

// TODO: 1 define methods
type QueryClient interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, errors.Error)
	RollbackTx(ctx context.Context, tx *sqlx.Tx) errors.Error
	CommitTx(ctx context.Context, tx *sqlx.Tx) errors.Error

	GetOne(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error
	GetOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any) errors.Error
	Get(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error
	GetTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any) errors.Error

	InsertOne(ctx context.Context, tagName string, args map[string]any) errors.Error
	InsertOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) errors.Error

	Delete(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error)
	DeleteTx(ctx context.Context, tx *sql.Tx, tagName string, args map[string]any) (int64, errors.Error)

	Update(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) (int64, errors.Error)

	GetRawQuery(tagName string, enum enum.QueryEnum) (*string, errors.Error)
}
