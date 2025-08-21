package originatingbiid

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v4/logging"
)

const ORIGINATING_BI_ID_CONTEXT_NAME = "originating-bi-id"

var originatingBiIdHeaderInDifferentCases = []string{"originating-bi-id", "Originating-Bi-Id", "originating-Bi-Id"}

type OriginatingBiIdProvider struct {
}

var logger logging.Logger

func init() {
	logger = logging.GetLogger(ORIGINATING_BI_ID_CONTEXT_NAME)
}

func (originatingBiIdProvider OriginatingBiIdProvider) InitLevel() int {
	return 0
}

func (originatingBiIdProvider OriginatingBiIdProvider) ContextName() string {
	return ORIGINATING_BI_ID_CONTEXT_NAME
}

func (originatingBiIdProvider OriginatingBiIdProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	for _, contextName := range originatingBiIdHeaderInDifferentCases {
		if incomingData[contextName] != nil {
			logger.Debug("context object=" + ORIGINATING_BI_ID_HEADER_NAME + " provided to context.Context")
			return context.WithValue(ctx, ORIGINATING_BI_ID_HEADER_NAME, NewOriginatingBiIdContextObject(incomingData[contextName].(string)))

		}
	}
	return ctx
}

func (originatingBiIdProvider OriginatingBiIdProvider) Set(ctx context.Context, originatingBiIdObject interface{}) (context.Context, error) {
	originatingBiId, success := originatingBiIdObject.(*originatingBiIdContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set originatingBiId")
	}
	logger.Debug("context object=" + ORIGINATING_BI_ID_CONTEXT_NAME + " set to context.Context")
	return context.WithValue(ctx, ORIGINATING_BI_ID_CONTEXT_NAME, originatingBiId), nil
}

func (originatingBiIdProvider OriginatingBiIdProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(ORIGINATING_BI_ID_CONTEXT_NAME)
}
