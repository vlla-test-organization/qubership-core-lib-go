package apiversion

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/ctxmanager"
	"testing"
)

const url_header_value = "api/v2/some-test-url"

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{ApiVersionProvider{}})
}

func TestInitApiVersionContext(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	apiVersion, _ := Of(ctx)
	assert.NotNil(t, apiVersion)
	assert.Equal(t, "v2", apiVersion.apiVersion)
}

func TestGetDefaultValueInitApiVersionContext(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), map[string]interface{}{})
	apiVersion, _ := Of(ctx)
	assert.NotNil(t, apiVersion)
	assert.Equal(t, "v"+default_version, apiVersion.apiVersion)
}

func TestSetApiVersionProvider(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	apiVersion, _ := Of(ctx)
	assert.Equal(t, "v2", apiVersion.apiVersion)

	var err error
	ctx, err = ctxmanager.SetContextObject(ctx, API_VERSION_CONTEXT_NAME, NewApiVersionContextObject("/v3/"))
	assert.Nil(t, err)
	secondApiVersion, _ := Of(ctx)
	assert.Equal(t, "v3", secondApiVersion.apiVersion)
}

func TestErrorSetAcceptLanguageProvider(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(API_VERSION_CONTEXT_NAME)
	_, err := provider.Set(context.Background(), "wrong type")
	assert.NotNil(t, err)
}

func TestContextName(t *testing.T) {
	assert.Equal(t, ApiVersionProvider{}.ContextName(), API_VERSION_CONTEXT_NAME)
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{url_header: url_header_value}
}

func TestGetLogValue(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	apiVersion, _ := Of(ctx)
	assert.Equal(t, "v2", apiVersion.GetVersion())
	assert.Equal(t, "v2", apiVersion.GetLogValue())
}
