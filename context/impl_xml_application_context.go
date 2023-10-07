package context

import (
	"github.com/aivyss/sql-mapper/errors"
	"github.com/aivyss/sql-mapper/reader/xml"
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

func (ctx *xmlApplicationContext) GetTxManager() TxManager {
	return GetApplicationContext().GetTxManager()
}

func (ctx *xmlApplicationContext) GetReadOnlyQueryClient(identifier string) (ReadOnlyQueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[getQueryClientKey(identifier, true)]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}

func (ctx *xmlApplicationContext) GetQueryClient(identifier string) (QueryClient, errors.Error) {
	client, ok := ctx.queryClientMap[getQueryClientKey(identifier, false)]
	if !ok {
		return nil, errors.BuildBasicErr(errors.NotFoundQueryClientErr)
	}

	return client, nil
}

func (ctx *xmlApplicationContext) RegisterQueryClient(client QueryClient) errors.Error {
	_, ok := ctx.queryClientMap[getQueryClientKey(client.Id(), client.ReadOnly())]
	if ok {
		return errors.BuildBasicErr(errors.RegisterQueryClientErr)
	}

	ctx.queryClientMap[getQueryClientKey(client.Id(), client.ReadOnly())] = client

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

			var queryClient QueryClient
			if readOnly {
				onlyQueryClient, err := NewReadOnlyQueryClient(client.Identifier, client.FilePath)

				if err != nil {
					resultErr = err
				} else {
					queryClient = onlyQueryClient.(QueryClient)
				}

			} else {
				queryClient, err = NewQueryClient(client.Identifier, client.FilePath)
				if err != nil {
					resultErr = err
				}
			}

			if err != nil {
				resultErr = err
			} else if queryClient == nil {
				resultErr = errors.BuildBasicErr(errors.RegisterQueryClientErr)
			} else {
				clientMap[queryClientKey{
					identifier: queryClient.Id(),
					readOnly:   readOnly,
				}] = queryClient
			}
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
