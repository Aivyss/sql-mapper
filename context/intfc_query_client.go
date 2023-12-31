package context

import (
	"context"
	"github.com/aivyss/sql-mapper/entity"
	"github.com/aivyss/sql-mapper/errors"
	"github.com/jmoiron/sqlx"
)

type QueryClient interface {
	ReadOnlyQueryClient
	InsertOne(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
	InsertOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) errors.Error

	Delete(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)
	DeleteTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)

	Update(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)
}

type ReadOnlyQueryClient interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, errors.Error)
	RollbackTx(ctx context.Context, tx *sqlx.Tx) errors.Error
	CommitTx(ctx context.Context, tx *sqlx.Tx) errors.Error

	GetOne(ctx context.Context, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
	GetOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
	Get(ctx context.Context, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
	GetTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error

	Id() string
	ReadOnly() bool
}
