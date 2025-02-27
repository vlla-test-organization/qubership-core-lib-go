package allowedheaders

import (
	"context"
	"errors"
	"strings"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	"github.com/netcracker/qubership-core-lib-go/v3/logging"
)

const (
	HEADERS_PROPERTY           = "headers.allowed"
	ALLOWED_HEADER_CONTEX_NAME = "allowed_header"
)

var logger logging.Logger

func init() {
	logger = logging.GetLogger("allowed-headers")
}

type allowedHeaderProvider struct {
	allowedHeaders []string
}

func NewAllowedHeaderProvider() allowedHeaderProvider {
	return allowedHeaderProvider{allowedHeaders: strings.Split(strings.ReplaceAll(configloader.GetOrDefaultString(HEADERS_PROPERTY, ""), " ", ""), ",")}
}

func (allowedHeaderProvider allowedHeaderProvider) InitLevel() int {
	return 0
}

func (allowedHeaderProvider allowedHeaderProvider) ContextName() string {
	return ALLOWED_HEADER_CONTEX_NAME
}

func (allowedHeaderProvider allowedHeaderProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	var allowedRequestHeaders = make(map[string]string)
	for _, headerName := range allowedHeaderProvider.allowedHeaders {
		if incomingData[headerName] != nil {
			allowedRequestHeaders[headerName] = incomingData[headerName].(string)
			logger.Debug("context object=" + headerName + " provided to context.Context")
		}
	}
	return context.WithValue(ctx, ALLOWED_HEADER_CONTEX_NAME, NewAllowedHeaderContextObject(allowedRequestHeaders))
}

func (allowedHeaderProvider allowedHeaderProvider) Set(ctx context.Context, allowedHeadersObject interface{}) (context.Context, error) {
	allowedHeaders, success := allowedHeadersObject.(allowedHeaderContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set allowedHeaders")
	}
	var allowedRequestHeaders = make(map[string]string)
	for _, headerName := range allowedHeaderProvider.allowedHeaders {
		if allowedHeaders.header[headerName] != "" {
			allowedRequestHeaders[headerName] = allowedHeaders.header[headerName]
			logger.Debug("context object=" + headerName + " set to context.Context")
		}
	}
	return context.WithValue(ctx, ALLOWED_HEADER_CONTEX_NAME, NewAllowedHeaderContextObject(allowedRequestHeaders)), nil
}

func (allowedHeaderProvider allowedHeaderProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(ALLOWED_HEADER_CONTEX_NAME)
}
