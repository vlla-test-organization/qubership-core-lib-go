package originatingbiid

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v7/context-propagation/ctxmanager"
)

const ORIGINATING_BI_ID_HEADER_NAME = "originating-bi-id"

type originatingBiIdContextObject struct {
	originatingBiId string
}

func NewOriginatingBiIdContextObject(originatingBiId string) *originatingBiIdContextObject {
	return &originatingBiIdContextObject{originatingBiId: originatingBiId}
}

func (originatingBiIdContextObject originatingBiIdContextObject) GetOriginatingBiId() string {
	return originatingBiIdContextObject.originatingBiId
}

func (originatingBiIdContextObject originatingBiIdContextObject) GetLogValue() string {
	return originatingBiIdContextObject.originatingBiId
}

func (originatingBiIdContextObject originatingBiIdContextObject) Serialize() (map[string]string, error) {
	if originatingBiIdContextObject.originatingBiId == "" {
		return nil, nil
	}
	return map[string]string{
		ORIGINATING_BI_ID_HEADER_NAME: originatingBiIdContextObject.originatingBiId,
	}, nil
}

func Of(ctx context.Context) (*originatingBiIdContextObject, error) {
	contextProvider, err := ctxmanager.GetProvider(ORIGINATING_BI_ID_CONTEXT_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	if abstractContextObject == nil {
		return nil, errors.New("originatingBiId context object is null")
	}
	contextObject := (abstractContextObject).(*originatingBiIdContextObject)
	return contextObject, nil
}
