package xversion

import (
	"context"
	"errors"

	"github.com/vlla-test-organization/qubership-core-lib-go/v4/context-propagation/ctxmanager"
)

const X_VERSION_HEADER_NAME = "X-Version"

type xVersionContextObject struct {
	xVersion string
}

func NewXVersionContextObject(headerValues string) xVersionContextObject {
	return xVersionContextObject{xVersion: headerValues}
}

func (xVersionContextObject xVersionContextObject) Serialize() (map[string]string, error) {
	if xVersionContextObject.xVersion == "" {
		return nil, nil
	}
	return map[string]string{X_VERSION_HEADER_NAME: xVersionContextObject.xVersion}, nil
}

func (xVersionContextObject xVersionContextObject) GetXVersion() string {
	return xVersionContextObject.xVersion
}

func (xVersionContextObject xVersionContextObject) GetLogValue() string {
	return xVersionContextObject.xVersion
}

func Of(ctx context.Context) (*xVersionContextObject, error) {
	contextProvider, err := ctxmanager.GetProvider(X_VERSION_CONTEXT_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	if err != nil {
		return nil, err
	}

	if abstractContextObject == nil {
		return nil, errors.New("X-Version context object is null")
	}

	contextObject := (abstractContextObject).(xVersionContextObject)
	return &contextObject, nil
}
