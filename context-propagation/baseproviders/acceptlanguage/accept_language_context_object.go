package acceptlanguage

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/context-propagation/ctxmanager"
	"strings"
)

const ACCEPT_LANGUAGE_HEADER_NAME = "Accept-Language"

type acceptLanguageContextObject struct {
	acceptLanguage []string
}

func NewAcceptLanguageContextObject(headerValue string) acceptLanguageContextObject {
	acceptLanguageList := strings.Split(headerValue, ",")
	return acceptLanguageContextObject{acceptLanguage: acceptLanguageList}
}

func (acceptLanguageContextObject acceptLanguageContextObject) GetAcceptLanguage() string {
	return strings.Join(acceptLanguageContextObject.acceptLanguage, ",")
}

func (acceptLanguageContextObject acceptLanguageContextObject) GetLogValue() string {
	return acceptLanguageContextObject.GetAcceptLanguage()
}

func (acceptLanguageContextObject acceptLanguageContextObject) Serialize() (map[string]string, error) {
	if acceptLanguageContextObject.acceptLanguage == nil {
		return nil, nil
	}
	return map[string]string{ACCEPT_LANGUAGE_HEADER_NAME: strings.Join(acceptLanguageContextObject.acceptLanguage, ",")}, nil
}

func Of(ctx context.Context) (*acceptLanguageContextObject, error) {
	contextProvider, err := ctxmanager.GetProvider(ACCEPT_LANGUAGE_CONTEXT_NAME)
	if err != nil {
		return nil, err
	}
	abstractContextObject := contextProvider.Get(ctx)
	if abstractContextObject == nil {
		return nil, errors.New("acceptLanguage context object is null")
	}
	contextObject := abstractContextObject.(acceptLanguageContextObject)
	return &contextObject, nil
}
