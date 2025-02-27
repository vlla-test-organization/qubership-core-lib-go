package apiversion

import (
	"context"
	"errors"
	"github.com/netcracker/qubership-core-lib-go/v3/context-propagation/ctxmanager"
	"strings"
)

const default_version = "1"

type apiVersionContextObject struct {
	apiVersion string
}

func NewApiVersionContextObject(headersValues string) apiVersionContextObject {
	var urlVersion = default_version
	if headersValues != "" {
		urlVersion = strings.Split(strings.Split(headersValues, "/v")[1], "/")[0]
	}
	return apiVersionContextObject{apiVersion: "v" + urlVersion}
}

func (apiVersionContextObject apiVersionContextObject) GetLogValue() string {
	return apiVersionContextObject.apiVersion
}

func (apiVersionContextObject apiVersionContextObject) GetVersion() string {
	return apiVersionContextObject.apiVersion
}

func Of(ctx context.Context) (*apiVersionContextObject, error) {
	contextProvider, err := ctxmanager.GetProvider(API_VERSION_CONTEXT_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	if abstractContextObject == nil {
		return nil, errors.New("apiVersion context object is null")
	}
	contextObject := (abstractContextObject).(apiVersionContextObject)
	return &contextObject, nil
}
