package endpoint

import (
	"context"
	"sql-mapper/enum"
	"sql-mapper/errors"
)

// TODO: 1 define methods
type QueryClient interface {
	GetOneByTagName(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error
	GetByTagName(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error
	GetRawQuery(tagName string, enum enum.QueryEnum) (*string, errors.Error)
}
