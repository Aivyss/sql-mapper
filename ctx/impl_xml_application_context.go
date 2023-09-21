package ctx

import (
	"sql-mapper/endpoint"
	"sql-mapper/errors"
)

type xmlApplicationContext struct {
	queryClientMap map[string]endpoint.QueryClient
}

func (ctx *xmlApplicationContext) GetQueryClient(identifier string) (endpoint.QueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[identifier]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}
func (ctx *xmlApplicationContext) RegisterQueryClient(identifier string, client endpoint.QueryClient) errors.Error {
	_, ok := ctx.queryClientMap[identifier]
	if ok {
		return errors.BuildBasicErr(errors.RegisterQueryClientErr)
	}

	ctx.queryClientMap[identifier] = client

	return nil
}

func NewApplicationContext() ApplicationContext {

	return &xmlApplicationContext{
		queryClientMap: map[string]endpoint.QueryClient{},
	}
}
