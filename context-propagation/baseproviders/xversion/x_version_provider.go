package xversion

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v3/logging"
)

const X_VERSION_CONTEXT_NAME = "X-Version"

type XVersionProvider struct {
}

var logger logging.Logger

func init() {
	logger = logging.GetLogger("x-version")
}

func (xVersionProvider XVersionProvider) InitLevel() int {
	return 0
}

func (xVersionProvider XVersionProvider) ContextName() string {
	return X_VERSION_CONTEXT_NAME
}

func (xVersionProvider XVersionProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	if incomingData[X_VERSION_HEADER_NAME] == nil {
		return ctx
	}
	logger.Debug("context object=" + X_VERSION_HEADER_NAME + " provided to context.Context")
	return context.WithValue(ctx, X_VERSION_CONTEXT_NAME, NewXVersionContextObject(incomingData[X_VERSION_HEADER_NAME].(string)))
}

func (xVersionProvider XVersionProvider) Set(ctx context.Context, xVersionObject interface{}) (context.Context, error) {
	xVersion, success := xVersionObject.(xVersionContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set xRequestId")
	}
	logger.Debug("context object=" + X_VERSION_CONTEXT_NAME + " set to context.Context")
	return context.WithValue(ctx, X_VERSION_CONTEXT_NAME, xVersion), nil
}

func (xVersionProvider XVersionProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(X_VERSION_CONTEXT_NAME)
}
