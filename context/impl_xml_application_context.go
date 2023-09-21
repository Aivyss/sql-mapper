package context

import (
	"github.com/jmoiron/sqlx"
	"sql-mapper/endpoint"
	"sql-mapper/errors"
	"sql-mapper/reader/xml"
	"sync"
)

var xmlAppCtxOnce sync.Once

// xmlAppCtx managed by singleton
var xmlAppCtx *xmlApplicationContext

// xmlApplicationContext application context based on xml file
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
func (ctx *xmlApplicationContext) RegisterQueryClient(client endpoint.QueryClient) errors.Error {
	_, ok := ctx.queryClientMap[client.Id()]
	if ok {
		return errors.BuildBasicErr(errors.RegisterQueryClientErr)
	}

	ctx.queryClientMap[client.Id()] = client

	return nil
}

func BuildXmlApplicationContext(db *sqlx.DB, filePath string) (ApplicationContext, errors.Error) {
	var resultErr errors.Error

	xmlAppCtxOnce.Do(func() {
		appCtxComp, err := xml.ReadSettings(&filePath)
		if err != nil {
			resultErr = err
		}

		clientMap := map[string]endpoint.QueryClient{}
		for _, client := range appCtxComp.QueryClientComponent.Clients {
			queryClient, err := NewQueryClient(db, client.Identifier, client.FilePath)
			if err != nil {
				resultErr = err
			}

			clientMap[queryClient.Id()] = queryClient
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
