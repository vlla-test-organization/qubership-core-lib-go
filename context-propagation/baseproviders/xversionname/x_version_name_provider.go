package xversionname

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v6/logging"
)

const X_VERSION_NAME_CONTEXT_NAME = "X-Version-Name"

type XVersionNameProvider struct {
}

var logger logging.Logger

func init() {
	logger = logging.GetLogger("x-version-name")
}

func (xVersionNameProvider XVersionNameProvider) InitLevel() int {
	return 0
}

func (xVersionNameProvider XVersionNameProvider) ContextName() string {
	return X_VERSION_NAME_CONTEXT_NAME
}

func (xVersionNameProvider XVersionNameProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	if incomingData[X_VERSION_NAME_HEADER_NAME] == nil {
		return ctx
	}
	logger.Debug("context object=" + X_VERSION_NAME_HEADER_NAME + " provided to context.Context")
	return context.WithValue(ctx, X_VERSION_NAME_CONTEXT_NAME, NewXVersionNameContextObject(incomingData[X_VERSION_NAME_HEADER_NAME].(string)))
}

func (xVersionNameProvider XVersionNameProvider) Set(ctx context.Context, xVersionObject interface{}) (context.Context, error) {
	xVersion, success := xVersionObject.(xVersionNameContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set xVersionName")
	}
	logger.Debug("context object=" + X_VERSION_NAME_CONTEXT_NAME + " set to context.Context")
	return context.WithValue(ctx, X_VERSION_NAME_CONTEXT_NAME, xVersion), nil
}

func (xVersionProvider XVersionNameProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(X_VERSION_NAME_CONTEXT_NAME)
}
