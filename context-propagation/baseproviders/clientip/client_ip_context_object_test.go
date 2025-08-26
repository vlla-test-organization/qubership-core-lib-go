package clientip

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v7/context-propagation/ctxmanager"
	"testing"
)

const client_ip_value = "127.0.0.1"
const x_forwarded_for_value = "127.0.0.2,127.0.0.3"
const x_forwarded_for_first_ip = "127.0.0.2"

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{ClientIpProvider{}})
}

func TestClientIpPropagatableCtx_Empty(t *testing.T) {
	incomingHeaders := map[string]interface{}{}
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_NC_CLIENT_IP_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	clientIp, _ := Of(ctx)
	assert.Equal(t, "", clientIp.clientIp)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	_, ok := outgoingData[X_NC_CLIENT_IP_HEADER_NAME]
	assert.Equal(t, false, ok)
}

func TestClientIpPropagatableCtx_OnlyNcClientIp(t *testing.T) {
	incomingHeaders := getOnlyNcClientIpHeader()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_NC_CLIENT_IP_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	clientIp, _ := Of(ctx)
	assert.Equal(t, client_ip_value, clientIp.clientIp)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	assert.Equal(t, client_ip_value, outgoingData[X_NC_CLIENT_IP_HEADER_NAME])
}

func TestClientIpPropagatableCtx_OnlyXForwardedFor(t *testing.T) {
	incomingHeaders := getOnlyXForwardedForHeader()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_NC_CLIENT_IP_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	clientIp, _ := Of(ctx)
	assert.Equal(t, x_forwarded_for_first_ip, clientIp.clientIp)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	assert.Equal(t, x_forwarded_for_first_ip, outgoingData[X_NC_CLIENT_IP_HEADER_NAME])
}

func TestClientIpPropagatableCtx_AllTwoHeaders(t *testing.T) {
	incomingHeaders := getAllTwoHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, X_NC_CLIENT_IP_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	clientIp, _ := Of(ctx)
	assert.Equal(t, x_forwarded_for_first_ip, clientIp.clientIp)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	assert.Equal(t, x_forwarded_for_first_ip, outgoingData[X_NC_CLIENT_IP_HEADER_NAME])
}

func TestOfClientIpContext(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getOnlyNcClientIpHeader())
	clientIp, _ := Of(ctx)
	assert.Equal(t, client_ip_value, clientIp.clientIp)
}

func TestSetClientIpProvider(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getOnlyXForwardedForHeader())

	clientIp, _ := Of(ctx)
	assert.Equal(t, x_forwarded_for_first_ip, clientIp.clientIp)

	var err error
	ctx, err = ctxmanager.SetContextObject(ctx, X_NC_CLIENT_IP_CONTEXT_NAME, NewClientIpContextObject("1.2.3.4"))
	assert.Nil(t, err)
	secondClientIp, _ := Of(ctx)
	assert.Equal(t, "1.2.3.4", secondClientIp.clientIp)
}

func TestWrongContextObjectType(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(X_NC_CLIENT_IP_CONTEXT_NAME)
	_, err := provider.Set(context.Background(), "wrong type")
	assert.NotNil(t, err)
}

func TestGetLogValue(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getAllTwoHeaders())

	xVersionName, _ := Of(ctx)
	assert.Equal(t, x_forwarded_for_first_ip, xVersionName.clientIp)

	assert.Equal(t, x_forwarded_for_first_ip, xVersionName.GetClientIp())
	assert.Equal(t, x_forwarded_for_first_ip, xVersionName.GetLogValue())
}

func TestContextName(t *testing.T) {
	assert.Equal(t, ClientIpProvider{}.ContextName(), X_NC_CLIENT_IP_CONTEXT_NAME)
}

func getOnlyNcClientIpHeader() map[string]interface{} {
	return map[string]interface{}{X_NC_CLIENT_IP_HEADER_NAME: client_ip_value}
}

func getOnlyXForwardedForHeader() map[string]interface{} {
	return map[string]interface{}{X_FORWARDED_FOR_HEADER_NAME: x_forwarded_for_value}
}

func getAllTwoHeaders() map[string]interface{} {
	return map[string]interface{}{X_FORWARDED_FOR_HEADER_NAME: x_forwarded_for_value}
}
