package context

import (
	"sql-mapper/endpoint"
	"sql-mapper/errors"
)

// directApplicationContext application context
type directApplicationContext struct {
	queryClientMap map[string]endpoint.QueryClient
}

func (ctx *directApplicationContext) GetReadOnlyQueryClient(identifier string) (endpoint.ReadOnlyQueryClient, errors.Error) {
	return ctx.GetQueryClient(identifier)
}

func (ctx *directApplicationContext) GetQueryClient(identifier string) (endpoint.QueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[identifier]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}
func (ctx *directApplicationContext) RegisterQueryClient(client endpoint.QueryClient) errors.Error {
	_, ok := ctx.queryClientMap[client.Id()]
	if ok {
		return errors.BuildBasicErr(errors.RegisterQueryClientErr)
	}

	ctx.queryClientMap[client.Id()] = client

	return nil
}
