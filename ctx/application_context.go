package ctx

import (
	"sql-mapper/endpoint"
	"sql-mapper/errors"
)

type ApplicationContext interface {
	GetQueryClient(identifier string) (endpoint.QueryClient, errors.Error)
	RegisterQueryClient(identifier string, client endpoint.QueryClient) errors.Error
}
