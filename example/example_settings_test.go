package example

import (
	"github.com/stretchr/testify/assert"
	"sql-mapper/context"
	"testing"
)

func TestXmlAppCtx(t *testing.T) {
	ctx, err := context.BuildXmlApplicationContext(db, "./setting/settings.xml")
	assert.Nil(t, err)

	identifiers := []string{"identifier4-1", "identifier4-2", "identifier4-3"}
	for _, identifier := range identifiers {
		client, err := ctx.GetQueryClient(identifier)
		assert.Nil(t, err)
		assert.NotNil(t, client)
	}
}

func TestAppCtx(t *testing.T) {
	xmlAppCtx, err := context.BuildXmlApplicationContext(db, "./setting/settings.xml")
	assert.Nil(t, err)

	appCtx := context.GetApplicationContext()

	assert.Equal(t, appCtx, xmlAppCtx)
}
