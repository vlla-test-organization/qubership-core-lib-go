package tenant

import (
	"context"
	"errors"
	"github.com/netcracker/qubership-core-lib-go/v3/logging"
)

const (
	tenantIdClaim       = "tenant-id"
	noTenant            = ""
)

type TenantProvider struct {
}

var (
	logger         logging.Logger
	noTenantObject = NewTenantContextObject(noTenant)
)

func init() {
	logger = logging.GetLogger("tenant")
}

func (tenantProvider TenantProvider) InitLevel() int {
	return TenantContextLevel
}

func (tenantProvider TenantProvider) ContextName() string {
	return TenantContextName
}

func (tenantProvider TenantProvider) Provide(ctx context.Context, incomingData map[string]interface{}) context.Context {
	var tenantFromHeader string
	headerValue := incomingData[TenantHeader]
	if headerValue == nil {
		tenantFromHeader = ""
	} else {
		tenantFromHeader = headerValue.(string)
	}
	logger.DebugC(ctx, "Discovered anonymous request, will be used tenant = %s", tenantFromHeader)
	return context.WithValue(ctx, TenantContextName, NewTenantContextObject(tenantFromHeader))
}

func (tenantProvider TenantProvider) Set(ctx context.Context, object interface{}) (context.Context, error) {
	objectString, success := object.(string)
	if !success {
		return ctx, errors.New("incorrect type to set to tenant context")
	}
	return context.WithValue(ctx, TenantContextName, NewTenantContextObject(objectString)), nil
}

func (tenantProvider TenantProvider) Get(ctx context.Context) interface{} {
	return ctx.Value(TenantContextName)
}


