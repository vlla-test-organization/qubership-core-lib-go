package clientip

import (
	"context"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/ctxmanager"
)

const X_NC_CLIENT_IP_HEADER_NAME = "X-Nc-Client-Ip"
const X_FORWARDED_FOR_HEADER_NAME = "X-Forwarded-For"

type clientIpContextObject struct {
	clientIp string
}

func NewClientIpContextObject(headerValues string) clientIpContextObject {
	return clientIpContextObject{clientIp: headerValues}
}

func (clientIpContextObject clientIpContextObject) Serialize() (map[string]string, error) {
	if clientIpContextObject.clientIp == "" {
		return nil, nil
	}
	return map[string]string{X_NC_CLIENT_IP_HEADER_NAME: clientIpContextObject.clientIp}, nil
}

func (clientIpContextObject clientIpContextObject) GetClientIp() string {
	return clientIpContextObject.clientIp
}

func (clientIpContextObject clientIpContextObject) GetLogValue() string {
	return clientIpContextObject.clientIp
}

func Of(ctx context.Context) (*clientIpContextObject, error) {
	contextProvider, err := ctxmanager.GetProvider(X_NC_CLIENT_IP_CONTEXT_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	contextObject := (abstractContextObject).(clientIpContextObject)
	return &contextObject, nil
}
