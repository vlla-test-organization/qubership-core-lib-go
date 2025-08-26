package ctxhelper

import (
	"context"

	"github.com/vlla-test-organization/qubership-core-lib-go/v7/context-propagation/ctxmanager"
	"github.com/vlla-test-organization/qubership-core-lib-go/v7/logging"
)

var logger logging.Logger

func init() {
	logger = logging.GetLogger("accept-language")
}

func AddSerializableContextData(ctx context.Context, f func(string, string)) error {
	logger.Debug("start serialize context data")
	contextData, err := ctxmanager.GetSerializableContextData(ctx)
	if err != nil {
		return err
	}
	addContextData(contextData, f)
	return nil
}

func AddResponsePropagatableContextData(ctx context.Context, f func(string, string)) error {
	logger.Debug("start collect and insert response propagatable context data")
	contextData, err := ctxmanager.GetResponsePropagatableContextData(ctx)
	if err != nil {
		return err
	}
	addContextData(contextData, f)
	return nil
}

func addContextData(contextData map[string]string, f func(string, string)) {
	for headerName, headerVals := range contextData {
		logger.Debug("add context=" + headerName + " with value=" + headerVals)
		f(headerName, headerVals)
	}
}
