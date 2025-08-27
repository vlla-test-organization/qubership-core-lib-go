package apiversion

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/logging"
)

const (
	url_header               = "cloud-core.context-propagation.url"
	API_VERSION_CONTEXT_NAME = "Api-Version-Context"
)

var logger logging.Logger

func init() {
	logger = logging.GetLogger("api-version")
}

type ApiVersionProvider struct {
}

func (apiVersionProvider ApiVersionProvider) InitLevel() int {
	return 0
}

func (apiVersionProvider ApiVersionProvider) ContextName() string {
	return API_VERSION_CONTEXT_NAME
}

func (apiVersionProvider ApiVersionProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	headerValue := ""
	if incomingData[url_header] != nil {
		headerValue = incomingData[url_header].(string)
	}
	logger.Debug("context object=" + API_VERSION_CONTEXT_NAME + " provided to context.Context")
	return context.WithValue(ctx, API_VERSION_CONTEXT_NAME, NewApiVersionContextObject(headerValue))
}

func (apiVersionProvider ApiVersionProvider) Set(ctx context.Context, apiVersionObject interface{}) (context.Context, error) {
	apiVersion, success := apiVersionObject.(apiVersionContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set apiVersion")
	}
	logger.Debug("context object=" + API_VERSION_CONTEXT_NAME + " set to context.Context")
	return context.WithValue(ctx, API_VERSION_CONTEXT_NAME, apiVersion), nil
}

func (apiVersionProvider ApiVersionProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(API_VERSION_CONTEXT_NAME)
}
