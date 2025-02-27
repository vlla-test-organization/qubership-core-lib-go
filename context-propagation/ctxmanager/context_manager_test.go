package ctxmanager

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testContextValue = "test_value"
const testContextValue3 = "test_value3"
const testContextName = "Test-Context"
const testContextName3 = "Test-Context-3"

func TestContextManagerRegisterSingle(t *testing.T) {
	testContextProvider := TestContextProvider{}
	RegisterSingle(testContextProvider)
	provider, err := GetProvider(testContextProvider.ContextName())
	assert.Equal(t, testContextProvider, provider)
	assert.Nil(t, err)
}

func TestContextManagerRegister(t *testing.T) {
	frameworkContexts := []ContextProvider{TestContextProvider{}}
	Register(frameworkContexts)
	for _, ctx := range frameworkContexts {
		provider, err := GetProvider(ctx.ContextName())
		assert.Equal(t, ctx, provider)
		assert.Nil(t, err)
	}
}

func TestContextManagerRegisterOverride(t *testing.T) {
	testContextProvider := TestContextProvider{}
	testContextProvider2 := TestContextProvider2{}
	assert.Equal(t, testContextProvider.ContextName(), testContextProvider2.ContextName())

	RegisterSingle(testContextProvider)
	provider, _ := GetProvider(testContextProvider.ContextName())
	assert.Equal(t, 0, provider.InitLevel())

	RegisterSingle(testContextProvider2)
	provider2, _ := GetProvider(testContextProvider2.ContextName())
	assert.Equal(t, 1, provider2.InitLevel())
}

func TestContextManagerGet(t *testing.T) {
	RegisterSingle(TestContextProvider{})

	requestHeaders := getIncomingRequestHeaders()

	var ctx = context.Background()
	ctx = InitContext(ctx, requestHeaders)

	contextObject, err := GetContextObject(ctx, testContextName)
	assert.Nil(t, err)
	assert.Equal(t, testContextValue,
		(contextObject).(TestContextObject).GetTestValue())
}

func TestContextManagerSet(t *testing.T) {

	RegisterSingle(TestContextProvider{})
	requestHeaders := getIncomingRequestHeaders()

	var ctx = context.Background()
	ctx = InitContext(ctx, requestHeaders)

	contextObject, err := GetContextObject(ctx, testContextName)
	assert.Nil(t, err)
	assert.Equal(t, testContextValue,
		(contextObject).(TestContextObject).GetTestValue())

	secondTestValue := "testValue2"
	ctx, err = SetContextObject(ctx, testContextName, NewTestContextObject(secondTestValue))
	assert.Nil(t, err)

	contextObject, err = GetContextObject(ctx, testContextName)
	assert.Nil(t, err)
	assert.Equal(t, secondTestValue,
		(contextObject).(TestContextObject).GetTestValue())
}

func TestGetSerializableContextData(t *testing.T) {
	RegisterSingle(TestContextProvider{})
	incomingHeaders := getIncomingRequestHeaders()
	ctx := InitContext(context.Background(), incomingHeaders)

	outgoingData, _ := GetSerializableContextData(ctx)
	assert.Equal(t, testContextValue, outgoingData[testContextName])
}

func TestGetDownstreamHeaders(t *testing.T) {
	RegisterSingle(TestContextProvider{})
	incomingHeaders := getIncomingRequestHeaders()
	ctx := InitContext(context.Background(), incomingHeaders)

	header, _ := GetDownstreamHeaders(ctx)
	assert.Equal(t, header[0], testContextName)
}

func TestContextManagerCreateFullContextSnapshot(t *testing.T) {

	Register([]ContextProvider{TestContextProvider{}, TestContextProvider3{}})

	incomingHeaders := map[string]interface{}{testContextName: testContextValue, testContextName3: testContextValue3}
	ctx := InitContext(context.Background(), incomingHeaders)

	fullContextMap := CreateFullContextSnapshot(ctx)

	assert.True(t, ctx.Value(testContextName) != nil)
	assert.True(t, fullContextMap[testContextName] != nil)
	assert.True(t, ctx.Value(testContextName3) != nil)
	assert.True(t, fullContextMap[testContextName3] != nil)

	for key, value := range fullContextMap {
		contextData, err := GetContextObject(ctx, key)
		assert.Nil(t, err)
		assert.Equal(t, value, contextData)
	}
}

func TestContextManagerCreateContextSnapshot(t *testing.T) {

	Register([]ContextProvider{TestContextProvider{}, TestContextProvider3{}})

	incomingHeaders := map[string]interface{}{testContextName: testContextValue, testContextName3: testContextValue3}
	ctx := InitContext(context.Background(), incomingHeaders)

	providerNames := []string{testContextName}

	contextMap := CreateContextSnapshot(ctx, providerNames)

	assert.True(t, ctx.Value(testContextName) != nil)
	assert.True(t, contextMap[testContextName] != nil)
	assert.True(t, ctx.Value(testContextName3) != nil)
	assert.True(t, contextMap[testContextName3] == nil)

	for key, value := range contextMap {
		contextData, err := GetContextObject(ctx, key)
		assert.Nil(t, err)
		assert.Equal(t, value, contextData)
	}
}

func TestContextManagerActivateContextSnapshot(t *testing.T) {

	RegisterSingle(TestContextProvider{})
	testProvider, _ := GetProvider(testContextName)

	incomingHeaders := getIncomingRequestHeaders()
	ctx := InitContext(context.Background(), incomingHeaders)

	ctxSnapshot, err := ActivateContextSnapshot(CreateFullContextSnapshot(ctx))
	ctxSnapshot, err = SetContextObject(ctxSnapshot, testContextName, NewTestContextObject("testValue2"))
	assert.Nil(t, err)

	originTestContext, _ := testProvider.Get(ctx).(TestContextObject)
	assert.Equal(t, testContextValue, originTestContext.GetTestValue())

	snapshotTestContext, _ := testProvider.Get(ctxSnapshot).(TestContextObject)
	assert.Equal(t, "testValue2", snapshotTestContext.GetTestValue())
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{testContextName: testContextValue}
}

// Test structures
//
// Test Context Object
type TestContextObject struct {
	value string
}

func NewTestContextObject(headerValue string) TestContextObject {
	return TestContextObject{value: headerValue}
}

func (testContextObject TestContextObject) Serialize() (map[string]string, error) {
	if testContextObject.value != "" {
		return map[string]string{testContextName: testContextObject.value}, nil
	} else {
		return nil, errors.New("testContext value is empty")
	}

}

func (testContextObject TestContextObject) GetTestValue() string {
	return testContextObject.value
}

// Test Provider
type TestContextProvider struct {
}

func (testContextProvider TestContextProvider) InitLevel() int {
	return 0
}

func (testContextProvider TestContextProvider) ContextName() string {
	return testContextName
}

func (testContextProvider TestContextProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	if incomingData[testContextName] == nil {
		return ctx
	}
	return context.WithValue(ctx, testContextName, NewTestContextObject(incomingData[testContextName].(string)))
}

func (testContextProvider TestContextProvider) Set(ctx context.Context, object interface{}) (context.Context, error) {
	obj, success := object.(TestContextObject)
	if !success {
		return nil, errors.New("incorrect type to set test context")
	}
	return context.WithValue(ctx, testContextName, obj), nil
}

func (testContextProvider TestContextProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(testContextName)
}

// /
// Test Provider 2 with different InitLevel
type TestContextProvider2 struct {
	TestContextProvider
}

func (testContextProvider TestContextProvider2) InitLevel() int {
	return 1
}

// Test Provider 2 with different InitLevel and ContextName
type TestContextProvider3 struct {
	TestContextProvider
}

func (testContextProvider TestContextProvider3) ContextName() string {
	return testContextName3
}

func (testContextProvider TestContextProvider3) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	if incomingData[testContextName3] == nil {
		return ctx
	}
	return context.WithValue(ctx, testContextName3, NewTestContextObject(incomingData[testContextName3].(string)))
}

func (testContextProvider TestContextProvider3) InitLevel() int {
	return 1
}

///
