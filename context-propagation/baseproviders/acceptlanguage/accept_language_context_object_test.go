package acceptlanguage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/context-propagation/ctxmanager"
	"testing"
)

const ACCEPT_LANGUAGE_VALUE = "en;ru;"

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{AcceptLanguageProvider{}})
}

func TestAcceptLanguagePropagatableCtx(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, ACCEPT_LANGUAGE_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	acceptLanguage, _ := Of(ctx)
	assert.Equal(t, ACCEPT_LANGUAGE_VALUE, acceptLanguage.acceptLanguage[0])
	acceptLanguage.acceptLanguage = []string{"asd"}
	acceptLanguage2, _ := Of(ctx)
	assert.NotEqual(t, acceptLanguage.acceptLanguage, acceptLanguage2.acceptLanguage)
	outgoingData, _ := ctxmanager.GetSerializableContextData(ctx)
	assert.Equal(t, ACCEPT_LANGUAGE_VALUE, outgoingData[ACCEPT_LANGUAGE_HEADER_NAME])
}

func TestSetAcceptLanguageCtx(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, ACCEPT_LANGUAGE_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)

	acceptLanguage, _ := Of(ctx)
	assert.Equal(t, ACCEPT_LANGUAGE_VALUE, acceptLanguage.acceptLanguage[0])
	ctx, err = ctxmanager.SetContextObject(ctx, ACCEPT_LANGUAGE_CONTEXT_NAME, NewAcceptLanguageContextObject("ru"))
	assert.Nil(t, err)

	acceptLanguage, _ = Of(ctx)
	assert.Equal(t, "ru", acceptLanguage.acceptLanguage[0])
}

func TestGetLogValue(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	contextData, err := ctxmanager.GetContextObject(ctx, ACCEPT_LANGUAGE_CONTEXT_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	acceptLanguage, _ := Of(ctx)
	assert.Equal(t, ACCEPT_LANGUAGE_VALUE, acceptLanguage.GetAcceptLanguage())
	assert.Equal(t, ACCEPT_LANGUAGE_VALUE, acceptLanguage.GetLogValue())
}

func TestErrorSetAcceptLanguageProvider(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(ACCEPT_LANGUAGE_CONTEXT_NAME)
	_, err := provider.Set(context.Background(), "wrong type")
	assert.NotNil(t, err)
}

func TestContextName(t *testing.T) {
	assert.Equal(t, AcceptLanguageProvider{}.ContextName(), ACCEPT_LANGUAGE_CONTEXT_NAME)
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{ACCEPT_LANGUAGE_CONTEXT_NAME: ACCEPT_LANGUAGE_VALUE}
}
