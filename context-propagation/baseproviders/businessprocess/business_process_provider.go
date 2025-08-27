package businessprocess

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/logging"
)

const BUSINESS_PROCESS_CONTEXT_NAME = "Business-Process-Id"

type BusinessProcessProvider struct {
}

var logger logging.Logger

func init() {
	logger = logging.GetLogger("business-process")
}

func (businessProcessProvider BusinessProcessProvider) InitLevel() int {
	return 0
}

func (businessProcessProvider BusinessProcessProvider) ContextName() string {
	return BUSINESS_PROCESS_CONTEXT_NAME
}

func (businessProcessProvider BusinessProcessProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	if incomingData[BUSINESS_PROCESS_HEADER_NAME] == nil {
		return context.WithValue(ctx, BUSINESS_PROCESS_HEADER_NAME, NewBusinessProcessContextObject(""))
	}
	logger.Debug("context object=" + BUSINESS_PROCESS_HEADER_NAME + " provided to context.Context")
	return context.WithValue(ctx, BUSINESS_PROCESS_HEADER_NAME, NewBusinessProcessContextObject(incomingData[BUSINESS_PROCESS_HEADER_NAME].(string)))
}

func (businessProcessProvider BusinessProcessProvider) Set(ctx context.Context, businessProcessObject interface{}) (context.Context, error) {
	businessProcess, success := businessProcessObject.(*businessProcessContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set businessProcess")
	}
	logger.Debug("context object=" + BUSINESS_PROCESS_CONTEXT_NAME + " set to context.Context")
	return context.WithValue(ctx, BUSINESS_PROCESS_CONTEXT_NAME, businessProcess), nil
}

func (businessProcessProvider BusinessProcessProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(BUSINESS_PROCESS_CONTEXT_NAME)
}
