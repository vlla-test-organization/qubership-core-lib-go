package allowedheaders

import (
	"context"
	"errors"
	"fmt"
	"github.com/netcracker/qubership-core-lib-go/v3/context-propagation/ctxmanager"
	"strings"
)

type allowedHeaderContextObject struct {
	header map[string]string
}

func NewAllowedHeaderContextObject(headers map[string]string) allowedHeaderContextObject {
	return allowedHeaderContextObject{header: headers}
}

func (allowedHeaderContextObject allowedHeaderContextObject) GetHeaders() map[string]string {
	return allowedHeaderContextObject.header
}

func (allowedHeaderContextObject allowedHeaderContextObject) GetLogValue() string {
	l := []string{}
	for key, value := range allowedHeaderContextObject.header {
		l = append(l, fmt.Sprintf("'%s':'%s'", key, value))
	}
	return strings.Join(l, ", ")
}

func (allowedHeaderContextObject allowedHeaderContextObject) Serialize() (map[string]string, error) {
	if allowedHeaderContextObject.header == nil || len(allowedHeaderContextObject.header) == 0 {
		return nil, nil
	}
	return allowedHeaderContextObject.header, nil
}

func Of(ctx context.Context) (*allowedHeaderContextObject, error) {
	var contextObject allowedHeaderContextObject
	contextProvider, err := ctxmanager.GetProvider(ALLOWED_HEADER_CONTEX_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	if abstractContextObject == nil {
		return nil, errors.New("allowedHeaders context object is null")
	}
	contextObject = (abstractContextObject).(allowedHeaderContextObject)

	return &contextObject, nil
}
