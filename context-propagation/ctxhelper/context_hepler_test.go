package ctxhelper

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/configloader"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/allowedheaders"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/ctxmanager"
)

const test_context_value = "test_value"
const custom_header = "Custom-header-1"
const custom_header_value = "custom_value"
const test_context = "Test-Context"

var ctx context.Context

func init() {

	_ = os.Setenv(allowedheaders.HEADERS_PROPERTY, custom_header)
	configloader.Init(configloader.EnvPropertySource())
	ctxmanager.Register([]ctxmanager.ContextProvider{allowedheaders.NewAllowedHeaderProvider(), TestContextProvider{}})

	ctx = context.Background()
	ctx = ctxmanager.InitContext(ctx, map[string]interface{}{custom_header: custom_header_value, test_context: test_context_value})

}

func TestAddSerializableContextData(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com", nil)
	assert.Nil(t, err)
	err = AddSerializableContextData(ctx, request.Header.Add)
	assert.Nil(t, err)
	assert.Equal(t, custom_header_value, request.Header.Get(custom_header))
	assert.Equal(t, test_context_value, request.Header.Get(test_context))
}

func TestAddResponsePropagatableContextData(t *testing.T) {
	response := http.Response{Header: http.Header{}}
	err := AddResponsePropagatableContextData(ctx, response.Header.Add)
	assert.Nil(t, err)
	assert.Equal(t, test_context_value, response.Header.Get(test_context))
}

type TestContextObject struct {
	value string
}

func NewTestContextObject(headerValue string) TestContextObject {
	return TestContextObject{value: headerValue}
}

func (testContextObject TestContextObject) Serialize() (map[string]string, error) {
	if testContextObject.value != "" {
		return map[string]string{test_context: testContextObject.value}, nil
	} else {
		return nil, errors.New("testContext value is empty")
	}

}

func (testContextObject TestContextObject) Propagate() (map[string]string, error) {
	return testContextObject.Serialize()

}

func (testContextObject TestContextObject) GetTestValue() string {
	return testContextObject.value
}

type TestContextProvider struct {
}

func (testContextProvider TestContextProvider) InitLevel() int {
	return 0
}

func (testContextProvider TestContextProvider) ContextName() string {
	return test_context
}

func (testContextProvider TestContextProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	if incomingData[test_context] == nil {
		return ctx
	}
	return context.WithValue(ctx, test_context, NewTestContextObject(incomingData[test_context].(string)))
}

func (testContextProvider TestContextProvider) Set(ctx context.Context, object interface{}) (context.Context, error) {
	obj, success := object.(TestContextObject)
	if !success {
		return nil, errors.New("incorrect type to set test context")
	}
	return context.WithValue(ctx, test_context, obj), nil
}

func (testContextProvider TestContextProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(test_context)
}
