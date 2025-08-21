package xversionname

import (
	"context"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/context-propagation/ctxmanager"
)

const X_VERSION_NAME_HEADER_NAME = "X-Version-Name"

type xVersionNameContextObject struct {
	xVersionName string
}

func NewXVersionNameContextObject(headerValues string) xVersionNameContextObject {
	return xVersionNameContextObject{xVersionName: headerValues}
}

func (xVersionNameContextObject xVersionNameContextObject) Serialize() (map[string]string, error) {
	if xVersionNameContextObject.xVersionName == "" {
		return nil, nil
	}
	return map[string]string{X_VERSION_NAME_HEADER_NAME: xVersionNameContextObject.xVersionName}, nil
}

func (xVersionNameContextObject xVersionNameContextObject) GetXVersionName() string {
	return xVersionNameContextObject.xVersionName
}

func (xVersionNameContextObject xVersionNameContextObject) GetLogValue() string {
	return xVersionNameContextObject.xVersionName
}

func Of(ctx context.Context) (*xVersionNameContextObject, error) {
	contextProvider, err := ctxmanager.GetProvider(X_VERSION_NAME_HEADER_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	if err != nil {
		return nil, err
	}
	contextObject := (abstractContextObject).(xVersionNameContextObject)
	return &contextObject, nil
}
