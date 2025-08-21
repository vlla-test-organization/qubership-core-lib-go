package xversion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v7/context-propagation/ctxmanager"
)

const x_version_value = "V2"

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{XVersionProvider{}})
}

func TestXVersionPropagatableCtx(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_VERSION_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	xVersion, _ := Of(ctx)
	assert.Equal(t, x_version_value, xVersion.xVersion)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	assert.Equal(t, x_version_value, outgoingData[X_VERSION_HEADER_NAME])
}

func TestOfXVersionContext(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())
	xVersion, _ := Of(ctx)
	assert.Equal(t, x_version_value, xVersion.xVersion)
}

func TestOfEmptyXVersionHeaderValue(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestEmptyHeaders())
	_, err := Of(ctx)
	assert.NotNil(t, err)

}

func TestSetXVersionProvider(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	xVersion, _ := Of(ctx)
	assert.Equal(t, x_version_value, xVersion.xVersion)

	var err error
	ctx, err = ctxmanager.SetContextObject(ctx, X_VERSION_CONTEXT_NAME, NewXVersionContextObject("V3"))
	assert.Nil(t, err)
	secondXVersion, _ := Of(ctx)
	assert.Equal(t, "V3", secondXVersion.xVersion)
}

func TestErrorSetAcceptLanguageProvider(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(X_VERSION_CONTEXT_NAME)
	_, err := provider.Set(context.Background(), "wrong type")
	assert.NotNil(t, err)
}

func TestContextName(t *testing.T) {
	assert.Equal(t, XVersionProvider{}.ContextName(), X_VERSION_CONTEXT_NAME)
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{X_VERSION_CONTEXT_NAME: x_version_value}
}

func getIncomingRequestEmptyHeaders() map[string]interface{} {
	return map[string]interface{}{X_VERSION_CONTEXT_NAME: nil}
}

func TestGetLogValue(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	xVersion, _ := Of(ctx)
	assert.Equal(t, x_version_value, xVersion.xVersion)

	assert.Equal(t, x_version_value, xVersion.GetXVersion())
	assert.Equal(t, x_version_value, xVersion.GetLogValue())
}
