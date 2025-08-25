package businessprocess

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v6/context-propagation/ctxmanager"
)

const BUSINESS_PROCESS_HEADER_NAME = "Business-Process-Id"

type businessProcessContextObject struct {
	businessProcessId string
}

func NewBusinessProcessContextObject(businessProcessId string) *businessProcessContextObject {
	return &businessProcessContextObject{businessProcessId: businessProcessId}
}

func (businessProcessContextObject businessProcessContextObject) GetBusinessProcessId() string {
	return businessProcessContextObject.businessProcessId
}

func (businessProcessContextObject *businessProcessContextObject) SetBusinessProcessId(businessProcessId string) {
	businessProcessContextObject.businessProcessId = businessProcessId
}

func (businessProcessContextObject *businessProcessContextObject) GetLogValue() string {
	return businessProcessContextObject.businessProcessId
}

func (businessProcessContextObject businessProcessContextObject) Serialize() (map[string]string, error) {
	if businessProcessContextObject.businessProcessId == "" {
		return nil, nil
	}
	return map[string]string{
		BUSINESS_PROCESS_HEADER_NAME: businessProcessContextObject.businessProcessId,
	}, nil
}

func Of(ctx context.Context) (*businessProcessContextObject, error) {
	contextProvider, err := ctxmanager.GetProvider(BUSINESS_PROCESS_CONTEXT_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	if abstractContextObject == nil {
		return nil, errors.New("businessProcess context object is null")
	}
	contextObject := (abstractContextObject).(*businessProcessContextObject)
	return contextObject, nil
}
