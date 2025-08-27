package acceptlanguage

import (
	"context"
	"errors"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/logging"
)

type AcceptLanguageProvider struct {
}

const ACCEPT_LANGUAGE_CONTEXT_NAME = "Accept-Language"

var logger logging.Logger

func init() {
	logger = logging.GetLogger("accept-language")
}

func (acceptLanguageProvider AcceptLanguageProvider) InitLevel() int {
	return 0
}

func (acceptLanguageProvider AcceptLanguageProvider) ContextName() string {
	return ACCEPT_LANGUAGE_CONTEXT_NAME
}

func (acceptLanguageProvider AcceptLanguageProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	if incomingData[ACCEPT_LANGUAGE_HEADER_NAME] == nil {
		return ctx
	}
	logger.Debug("context object=" + ACCEPT_LANGUAGE_HEADER_NAME + " provided to context.Context")
	acceptLanguageString, successToString := incomingData[ACCEPT_LANGUAGE_HEADER_NAME].(string)
	if !successToString {
		return ctx
	}
	return context.WithValue(ctx, ACCEPT_LANGUAGE_CONTEXT_NAME, NewAcceptLanguageContextObject(acceptLanguageString))
}

func (acceptLanguageProvider AcceptLanguageProvider) Set(ctx context.Context, acceptLanguageObject interface{}) (context.Context, error) {
	acceptLanguage, success := acceptLanguageObject.(acceptLanguageContextObject)
	if !success {
		return ctx, errors.New("incorrect type to set acceptLanguage")
	}
	logger.Debug("context object=" + ACCEPT_LANGUAGE_CONTEXT_NAME + " set to context.Context")
	return context.WithValue(ctx, ACCEPT_LANGUAGE_CONTEXT_NAME, acceptLanguage), nil
}

func (acceptLanguageProvider AcceptLanguageProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(ACCEPT_LANGUAGE_CONTEXT_NAME)
}
