package xrequestid

import (
	"context"
	"github.com/netcracker/qubership-core-lib-go/v3/context-propagation/ctxmanager"
	"github.com/stretchr/testify/assert"
	"testing"
)

const x_request_id_value = "123"

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{XRequestIdProvider{}})
}

func TestRequestIdSerializableCtx(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_REQUEST_ID_COTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	requestId, _ := Of(ctx)
	assert.Equal(t, x_request_id_value, requestId.requestId)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	assert.Equal(t, x_request_id_value, outgoingData[X_REQUEST_ID_HEADER_NAME])
}

func TestRequestIdIncomingResponsePropagatableCtx(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_REQUEST_ID_COTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	requestId, _ := Of(ctx)
	assert.Equal(t, x_request_id_value, requestId.requestId)
	responseContextData, _ := ctxmanager.GetResponsePropagatableContextData(ctx)
	assert.Equal(t, x_request_id_value, responseContextData[X_REQUEST_ID_HEADER_NAME])
}

func TestGetDefaultRequestId(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	requestId, _ := Of(ctx)
	assert.NotNil(t, requestId)
	assert.Equal(t, x_request_id_value, requestId.requestId)
}

func TestOfRequestIdContext(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())
	requestId, _ := Of(ctx)
	assert.Equal(t, x_request_id_value, requestId.requestId)
}

func TestSetRequestIdProvider(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	xRequestId, _ := Of(ctx)
	assert.Equal(t, x_request_id_value, xRequestId.requestId)

	var err error
	ctx, err = ctxmanager.SetContextObject(ctx, X_REQUEST_ID_COTEXT_NAME, NewXRequestIdContextObject("321"))
	assert.Nil(t, err)
	secondXRequestId, _ := Of(ctx)
	assert.Equal(t, "321", secondXRequestId.requestId)
}

func TestErrorSetAcceptLanguageProvider(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(X_REQUEST_ID_COTEXT_NAME)
	_, err := provider.Set(context.Background(), "wrong type")
	assert.NotNil(t, err)
}

func TestContextName(t *testing.T) {
	assert.Equal(t, XRequestIdProvider{}.ContextName(), X_REQUEST_ID_COTEXT_NAME)
}

func TestXRequestIdInterface(t *testing.T) {
	var xRequestIdInterface XRequestId
	xRequestIdInterface = xRequestIdContextObject{requestId: x_request_id_value}
	xRequestIdInterface.GetRequestId()
	assert.Equal(t, x_request_id_value, xRequestIdInterface.GetRequestId())
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{X_REQUEST_ID_COTEXT_NAME: x_request_id_value}
}

func TestGetLogValue(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	xRequestId, _ := Of(ctx)
	assert.Equal(t, x_request_id_value, xRequestId.requestId)

	assert.Equal(t, x_request_id_value, xRequestId.GetRequestId())
	assert.Equal(t, x_request_id_value, xRequestId.GetLogValue())
}
