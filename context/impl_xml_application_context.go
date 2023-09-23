package context

import (
	"sql-mapper/errors"
	"sql-mapper/reader/xml"
	"strconv"
)

type queryClientKey struct {
	identifier string
	readOnly   bool
}

// xmlApplicationContext application context based on xml file
type xmlApplicationContext struct {
	dBManager
	queryClientMap map[queryClientKey]QueryClient
}

func (ctx *xmlApplicationContext) GetReadOnlyQueryClient(identifier string) (ReadOnlyQueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[getQueryKey(identifier, true)]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}

func (ctx *xmlApplicationContext) GetQueryClient(identifier string) (QueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[getQueryKey(identifier, false)]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}

func (ctx *xmlApplicationContext) RegisterQueryClient(client QueryClient) errors.Error {
	_, ok := ctx.queryClientMap[getQueryKey(client.Id(), false)]
	if ok {
		return errors.BuildBasicErr(errors.RegisterQueryClientErr)
	}

	ctx.queryClientMap[getQueryKey(client.Id(), false)] = client

	return nil
}

func registerXmlContext(filePath string) (ApplicationContext, errors.Error) {
	var resultErr errors.Error

	xmlAppCtxOnce.Do(func() {
		appCtxComp, err := xml.ReadSettings(&filePath)
		if err != nil {
			resultErr = err
		}

		clientMap := map[queryClientKey]QueryClient{}
		for _, client := range appCtxComp.QueryClientComponent.Clients {
			readOnly, pErr := strconv.ParseBool(client.ReadOnly)
			if pErr != nil {
				readOnly = true
			}
			queryClient, err := NewQueryClient(client.Identifier, client.FilePath, readOnly)

			if err != nil {
				resultErr = err
			}

			clientMap[queryClientKey{
				identifier: queryClient.Id(),
				readOnly:   readOnly,
			}] = queryClient
		}

		xmlAppCtx = &xmlApplicationContext{
			queryClientMap: clientMap,
		}

		integAppCtx.xmlAppCtx = xmlAppCtx
	})

	if resultErr != nil {
		return nil, resultErr
	}

	return integAppCtx, nil
}
