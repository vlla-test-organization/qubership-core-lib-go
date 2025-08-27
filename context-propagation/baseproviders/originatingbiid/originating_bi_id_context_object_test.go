package originatingbiid

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/ctxmanager"
	"testing"
)

const (
	guid = "a97e7c33-2bad-447d-9199-dc49c8216ea0"
)

func init() {
	ctxmanager.Register([]ctxmanager.ContextProvider{OriginatingBiIdProvider{}})
}

func TestInitOriginatingBiIdContext(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	originatingBiIdObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, originatingBiIdObjectContext)
	assert.Equal(t, guid, originatingBiIdObjectContext.GetOriginatingBiId())
}

func TestInitOriginatingBiIdContextWithEmptyHeaders(t *testing.T) {
	incomingHeaders := getEmptyIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	originatingBiIdObjectContext, err := Of(ctx)
	assert.NotNil(t, err)
	assert.Nil(t, originatingBiIdObjectContext)
}

func TestInitOriginatingBiIdContextSerializer(t *testing.T) {
	incomingHeaders := getIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	originatingBiIdObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, originatingBiIdObjectContext)
	headers, err := originatingBiIdObjectContext.Serialize()
	assert.Empty(t, err)
	assert.NotEmpty(t, headers)
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, guid, headers[ORIGINATING_BI_ID_HEADER_NAME])

	outgoingData, err := ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.NotEmpty(t, outgoingData)
	assert.Equal(t, 1, len(outgoingData))
	assert.Equal(t, guid, outgoingData[ORIGINATING_BI_ID_HEADER_NAME])
}

func TestSetOriginatingBiIdDuringExecution(t *testing.T) {
	incomingHeaders := getEmptyIncomingRequestHeaders()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	originatingBiIdObjectContext, err := Of(ctx)
	assert.NotNil(t, err)
	assert.Nil(t, originatingBiIdObjectContext)

	outgoingData, err := ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.Empty(t, outgoingData)

	ctx, _ = ctxmanager.SetContextObject(ctx, ORIGINATING_BI_ID_CONTEXT_NAME, NewOriginatingBiIdContextObject(guid))

	outgoingData, err = ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.NotEmpty(t, outgoingData)
	assert.Equal(t, 1, len(outgoingData))
	assert.Equal(t, guid, outgoingData[ORIGINATING_BI_ID_HEADER_NAME])
}

func TestIncorrectProviderSetting(t *testing.T) {
	provider, _ := ctxmanager.GetProvider(ORIGINATING_BI_ID_CONTEXT_NAME)
	_, err := provider.Set(context.Background(), "incorrect type")
	assert.NotNil(t, err)
}

func getIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{ORIGINATING_BI_ID_HEADER_NAME: guid}
}

func getEmptyIncomingRequestHeaders() map[string]interface{} {
	return map[string]interface{}{}
}

func getIncomingRequestHeadersWithOtherCase() map[string]interface{} {
	return map[string]interface{}{"Originating-Bi-Id": guid}
}

func TestInitOriginatingBiIdContextWithOtherCase(t *testing.T) {
	incomingHeaders := getIncomingRequestHeadersWithOtherCase()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	originatingBiIdObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, originatingBiIdObjectContext)
	assert.Equal(t, guid, originatingBiIdObjectContext.GetOriginatingBiId())
}

func TestInitOriginatingBiIdContextSerializerWithOtherCase(t *testing.T) {
	incomingHeaders := getIncomingRequestHeadersWithOtherCase()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	originatingBiIdObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, originatingBiIdObjectContext)
	headers, err := originatingBiIdObjectContext.Serialize()
	assert.Empty(t, err)
	assert.NotEmpty(t, headers)
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, guid, headers[ORIGINATING_BI_ID_HEADER_NAME])

	outgoingData, err := ctxmanager.GetSerializableContextData(ctx)
	assert.Empty(t, err)
	assert.NotEmpty(t, outgoingData)
	assert.Equal(t, 1, len(outgoingData))
	assert.Equal(t, guid, outgoingData[ORIGINATING_BI_ID_HEADER_NAME])
}

func TestGetLogValue(t *testing.T) {
	incomingHeaders := getIncomingRequestHeadersWithOtherCase()
	ctx := ctxmanager.InitContext(context.Background(), incomingHeaders)
	originatingBiIdObjectContext, err := Of(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, originatingBiIdObjectContext)

	assert.Equal(t, guid, originatingBiIdObjectContext.GetOriginatingBiId())
	assert.Equal(t, guid, originatingBiIdObjectContext.GetLogValue())
}
