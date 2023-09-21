package context

import (
	"sql-mapper/endpoint"
	"sql-mapper/errors"
)

type KindAppCtx int

const (
	XML KindAppCtx = iota
)

type ApplicationContext interface {
	GetQueryClient(identifier string) (endpoint.QueryClient, errors.Error)
	RegisterQueryClient(client endpoint.QueryClient) errors.Error
}

// integAppCtx DO NOT REASSIGN THIS
var integAppCtx = &integratedApplicationContext{
	directAppCtx: &directApplicationContext{
		queryClientMap: map[string]endpoint.QueryClient{},
	},
}

// integratedApplicationContext (ApplicationContext by multi-metadata)
type integratedApplicationContext struct {
	xmlAppCtx    *xmlApplicationContext
	directAppCtx *directApplicationContext
}

func (c *integratedApplicationContext) GetQueryClient(identifier string) (endpoint.QueryClient, errors.Error) {
	contexts := []ApplicationContext{
		c.xmlAppCtx,
		c.directAppCtx,
	}
	var resultErr errors.Error

	for _, context := range contexts {
		client, err := context.GetQueryClient(identifier)
		if client != nil {
			return client, nil
		}

		resultErr = err
	}

	return nil, resultErr
}

func (c *integratedApplicationContext) RegisterQueryClient(client endpoint.QueryClient) errors.Error {
	return c.directAppCtx.RegisterQueryClient(client)
}

func GetApplicationContext() ApplicationContext {
	return integAppCtx
}
