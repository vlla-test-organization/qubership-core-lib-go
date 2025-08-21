package clientip

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/logging"
	"strings"
)

const X_NC_CLIENT_IP_CONTEXT_NAME = "X-Nc-Client-Ip"

type ClientIpProvider struct {
}

var logger logging.Logger

func init() {
	logger = logging.GetLogger("client-ip")
}

func (clientIpProvider ClientIpProvider) InitLevel() int {
	return 0
}

func (clientIpProvider ClientIpProvider) ContextName() string {
	return X_NC_CLIENT_IP_CONTEXT_NAME
}

func (clientIpProvider ClientIpProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	forwardedForHeaderValue := incomingData[X_FORWARDED_FOR_HEADER_NAME]
	clientIpHeaderValue := incomingData[X_NC_CLIENT_IP_HEADER_NAME]
	if forwardedForHeaderValue == nil && clientIpHeaderValue == nil {
		return context.WithValue(ctx, X_NC_CLIENT_IP_CONTEXT_NAME, NewClientIpContextObject(""))
	}
	logger.Debug("context object=" + X_NC_CLIENT_IP_HEADER_NAME + " provided to context.Context")
	if forwardedForHeaderValue != nil {
		firstIp := strings.Split(forwardedForHeaderValue.(string), ",")[0]
		return context.WithValue(ctx, X_NC_CLIENT_IP_CONTEXT_NAME, NewClientIpContextObject(firstIp))
	} else {
		return context.WithValue(ctx, X_NC_CLIENT_IP_CONTEXT_NAME, NewClientIpContextObject(clientIpHeaderValue.(string)))
	}
}

func (clientIpProvider ClientIpProvider) Set(ctx context.Context, clientIpObject interface{}) (context.Context, error) {
	clientIp, success := clientIpObject.(clientIpContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set clientIp")
	}
	logger.Debug("context object=" + X_NC_CLIENT_IP_CONTEXT_NAME + " set to context.Context")
	return context.WithValue(ctx, X_NC_CLIENT_IP_CONTEXT_NAME, clientIp), nil
}

func (clientIpProvider ClientIpProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(X_NC_CLIENT_IP_CONTEXT_NAME)
}
