package context

import (
	"sql-mapper/errors"
)

// directApplicationContext application context
type directApplicationContext struct {
	dBManager
	queryClientMap map[queryClientKey]QueryClient
}

func (ctx *directApplicationContext) GetReadOnlyQueryClient(identifier string) (ReadOnlyQueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[getQueryKey(identifier, true)]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}

func (ctx *directApplicationContext) GetQueryClient(identifier string) (QueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[getQueryKey(identifier, false)]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}
func (ctx *directApplicationContext) RegisterQueryClient(client QueryClient) errors.Error {
	qKey := getQueryKey(client.Id(), client.ReadOnly())
	_, ok := ctx.queryClientMap[qKey]
	if ok {
		return errors.BuildBasicErr(errors.RegisterQueryClientErr)
	}

	ctx.queryClientMap[qKey] = client

	return nil
}
