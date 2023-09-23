package context

import (
	"sql-mapper/entity"
	"sync"
)

// integAppCtx DO NOT REASSIGN THIS
var integAppCtx = &integratedApplicationContext{
	directAppCtx: &directApplicationContext{
		queryClientMap: map[queryClientKey]QueryClient{},
	},
	dbSet: &entity.DbSet{
		Write: nil,
		Read:  nil,
	},
}

// ---

// integAppCtx DO NOT REASSIGN THIS
var xmlAppCtxOnce sync.Once

// xmlAppCtx managed by singleton
var xmlAppCtx *xmlApplicationContext
