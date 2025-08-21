package xversionname

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v4/context-propagation/ctxmanager"
	"testing"
)

const x_version_name_value = "candidate"

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{XVersionNameProvider{}})
}

func TestXVersionNamePropagatableCtx(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_VERSION_NAME_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	xVersionName, _ := Of(ctx)
	assert.Equal(t, x_version_name_value, xVersionName.xVersionName)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	assert.Equal(t, x_version_name_value, outgoingData[X_VERSION_NAME_HEADER_NAME])
}

func TestOfXVersionNameContext(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())
	xVersionName, _ := Of(ctx)
	assert.Equal(t, x_version_name_value, xVersionName.xVersionName)
}

func TestSetXVersionNameProvider(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	xVersionName, _ := Of(ctx)
	assert.Equal(t, x_version_name_value, xVersionName.xVersionName)

	var err error
	ctx, err = ctxmanager.SetContextObject(ctx, X_VERSION_NAME_CONTEXT_NAME, NewXVersionNameContextObject("V3"))
	assert.Nil(t, err)
	secondXVersionName, _ := Of(ctx)
	assert.Equal(t, "V3", secondXVersionName.xVersionName)
}

func TestErrorSetAcceptLanguageProvider(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(X_VERSION_NAME_CONTEXT_NAME)
	_, err := provider.Set(context.Background(), "wrong type")
	assert.NotNil(t, err)
}

func TestContextName(t *testing.T) {
	assert.Equal(t, XVersionNameProvider{}.ContextName(), X_VERSION_NAME_CONTEXT_NAME)
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{X_VERSION_NAME_CONTEXT_NAME: x_version_name_value}
}

func TestGetLogValue(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	xVersionName, _ := Of(ctx)
	assert.Equal(t, x_version_name_value, xVersionName.xVersionName)

	assert.Equal(t, x_version_name_value, xVersionName.GetXVersionName())
	assert.Equal(t, x_version_name_value, xVersionName.GetLogValue())
}
