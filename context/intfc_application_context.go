package context

import (
	"github.com/aivyss/sql-mapper/entity"
	"github.com/aivyss/sql-mapper/errors"
	"github.com/jmoiron/sqlx"
)

type ApplicationContext interface {
	GetQueryClient(identifier string) (QueryClient, errors.Error)
	GetReadOnlyQueryClient(identifier string) (ReadOnlyQueryClient, errors.Error)
	RegisterQueryClient(client QueryClient) errors.Error
	GetDBs() *entity.DbSet
	GetDB(readDB bool) *sqlx.DB
	GetTxManager() TxManager
}

func GetApplicationContext() ApplicationContext {
	return integAppCtx
}

// integratedApplicationContext (ApplicationContext by multi-metadata)
type integratedApplicationContext struct {
	xmlAppCtx    *xmlApplicationContext
	directAppCtx *directApplicationContext
	dbSet        *entity.DbSet
	txManager    TxManager
}

func (c *integratedApplicationContext) GetTxManager() TxManager {
	return c.txManager
}

func (c *integratedApplicationContext) GetDBs() *entity.DbSet {
	return c.dbSet
}

func (c *integratedApplicationContext) GetDB(readDB bool) *sqlx.DB {
	if readDB {
		return c.dbSet.Read
	}

	return c.dbSet.Write
}

func (c *integratedApplicationContext) GetReadOnlyQueryClient(identifier string) (ReadOnlyQueryClient, errors.Error) {
	var resultErr errors.Error

	for _, context := range c.getContexts() {
		client, err := context.GetReadOnlyQueryClient(identifier)
		if client != nil {
			return client, nil
		}

		resultErr = err
	}

	return nil, resultErr
}

func (c *integratedApplicationContext) GetQueryClient(identifier string) (QueryClient, errors.Error) {
	var resultErr errors.Error

	for _, context := range c.getContexts() {
		client, err := context.GetQueryClient(identifier)
		if client != nil {
			return client, nil
		}

		resultErr = err
	}

	return nil, resultErr
}

func (c *integratedApplicationContext) RegisterQueryClient(client QueryClient) errors.Error {
	return c.directAppCtx.RegisterQueryClient(client)
}

func (c *integratedApplicationContext) getContexts() []ApplicationContext {
	return []ApplicationContext{
		c.xmlAppCtx,
		c.directAppCtx,
	}
}
