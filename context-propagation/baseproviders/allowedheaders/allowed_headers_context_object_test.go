package allowedheaders

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	"github.com/netcracker/qubership-core-lib-go/v3/context-propagation/ctxmanager"
	"github.com/stretchr/testify/assert"
)

const CUSTOM_HEADER_2 = "Custom-header-2"
const CUSTOM_HEADER = "Custom-header-1"
const CUSTOM_HEADER_VALUE = "custom_value"
const CUSTOM_HEADER_VALUE_2 = "custom_value_2"

func init() {
	_ = os.Setenv(HEADERS_PROPERTY, CUSTOM_HEADER+", "+CUSTOM_HEADER_2)
	configloader.Init(configloader.EnvPropertySource())
	ctxmanager.Register([]ctxmanager.ContextProvider{NewAllowedHeaderProvider()})
}

func TestAllowedHeaderCtx(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())
	contextData, err := ctxmanager.GetContextObject(ctx, ALLOWED_HEADER_CONTEX_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	allowedHeaders, _ := Of(ctx)

	assert.Equal(t, CUSTOM_HEADER_VALUE, allowedHeaders.header[CUSTOM_HEADER])
	assert.Equal(t, CUSTOM_HEADER_VALUE_2, allowedHeaders.header[CUSTOM_HEADER_2])
	os.Unsetenv("headers.allowed")

}

func TestPanicWrongAllowedHeaderCtx(t *testing.T) {
	//check panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())
	contextData, err := ctxmanager.GetContextObject(ctx, ALLOWED_HEADER_CONTEX_NAME)
	assert.NotNil(t, contextData)
	assert.Nil(t, err)
	allowedHeaders, _ := Of(ctx)
	assert.Empty(t, allowedHeaders.header["SomeContext"][0]) // here must be panic
}

func TestSetAllowedHeadersProvider(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	allowedHeaders, _ := Of(ctx)

	assert.Equal(t, CUSTOM_HEADER_VALUE, allowedHeaders.header[CUSTOM_HEADER])

	var err error
	ctx, err = ctxmanager.SetContextObject(ctx, ALLOWED_HEADER_CONTEX_NAME, NewAllowedHeaderContextObject(map[string]string{CUSTOM_HEADER: "new_custom_value"}))
	assert.Nil(t, err)
	secondAllowedHeaders, _ := Of(ctx)
	assert.Equal(t, "new_custom_value", secondAllowedHeaders.header[CUSTOM_HEADER])
	assert.Equal(t, CUSTOM_HEADER_VALUE_2, allowedHeaders.header[CUSTOM_HEADER_2])
}

func TestErrorSetAcceptLanguageProvider(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(ALLOWED_HEADER_CONTEX_NAME)
	_, err := provider.Set(context.Background(), "wrong type")
	assert.NotNil(t, err)
}

func TestContextName(t *testing.T) {
	assert.Equal(t, allowedHeaderProvider{}.ContextName(), ALLOWED_HEADER_CONTEX_NAME)
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{
		CUSTOM_HEADER:   CUSTOM_HEADER_VALUE,
		CUSTOM_HEADER_2: CUSTOM_HEADER_VALUE_2}
}

func TestGetLogValue(t *testing.T) {
	ctx := ctxmanager.InitContext(context.Background(), getIncomingRequestHeaders())

	allowedHeaders, _ := Of(ctx)

	expectedHeader1 := fmt.Sprintf("'%s':'%s'", CUSTOM_HEADER, CUSTOM_HEADER_VALUE)
	expectedHeader2 := fmt.Sprintf("'%s':'%s'", CUSTOM_HEADER_2, CUSTOM_HEADER_VALUE_2)

	actualHeaders := allowedHeaders.GetLogValue()
	assert.Contains(t, actualHeaders, expectedHeader1)
	assert.Contains(t, actualHeaders, expectedHeader2)
}
