package businessprocess

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/context-propagation/ctxmanager"
	"testing"
)

const (
	guid = "a97e7c33-2bad-447d-9199-dc49c8216ea0"
)

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{BusinessProcessProvider{}})
}

func TestInitBusinessProcessContext(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	businessProcessObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, businessProcessObjectContext)
	assert.Equal(t, guid, businessProcessObjectContext.GetBusinessProcessId())
}

func TestInitBusinessProcessContextWithEmptyHeaders(t *testing.T) {
	incomingHeaders := getEmptyIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	businessProcessObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, businessProcessObjectContext)
	assert.Empty(t, businessProcessObjectContext.GetBusinessProcessId())
}

func TestInitBusinessProcessContextSerializer(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	businessProcessObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, businessProcessObjectContext)
	headers, err := businessProcessObjectContext.Serialize()
	assert.Empty(t, err)
	assert.NotEmpty(t, headers)
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, guid, headers[BUSINESS_PROCESS_HEADER_NAME])

	outgoingData, err := ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.NotEmpty(t, outgoingData)
	assert.Equal(t, 1, len(outgoingData))
	assert.Equal(t, guid, outgoingData[BUSINESS_PROCESS_HEADER_NAME])
}

func TestInitBusinessProcessContextSerializerWithEmptyHeaders(t *testing.T) {
	incomingHeaders := getEmptyIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	businessProcessObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, businessProcessObjectContext)
	headers, err := businessProcessObjectContext.Serialize()
	assert.Empty(t, err)
	assert.Empty(t, headers)

	outgoingData, err := ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.Empty(t, outgoingData)
}

func TestSetBusinessProcessIdDuringExecution(t *testing.T) {
	incomingHeaders := getEmptyIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	businessProcessObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, businessProcessObjectContext)

	outgoingData, err := ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.Empty(t, outgoingData)

	businessProcessObjectContext.SetBusinessProcessId(guid)

	outgoingData, err = ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.NotEmpty(t, outgoingData)
	assert.Equal(t, 1, len(outgoingData))
	assert.Equal(t, guid, outgoingData[BUSINESS_PROCESS_HEADER_NAME])
}

func TestIncorrectProviderSetting(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(BUSINESS_PROCESS_CONTEXT_NAME)
	_, err := provider.Set(context.Background(), "incorrect type")
	assert.NotNil(t, err)
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{BUSINESS_PROCESS_HEADER_NAME: guid}
}

func getEmptyIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{}
}

func TestGetLogValue(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	businessProcessObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, businessProcessObjectContext)

	assert.Equal(t, guid, businessProcessObjectContext.GetBusinessProcessId())
	assert.Equal(t, guid, businessProcessObjectContext.GetLogValue())
}
