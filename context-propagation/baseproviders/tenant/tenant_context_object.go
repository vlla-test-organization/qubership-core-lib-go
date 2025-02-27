package tenant

import (
	"context"
	"github.com/netcracker/qubership-core-lib-go/v3/context-propagation/ctxmanager"
	"github.com/pkg/errors"
)

type TenantContextObject struct {
	tenant string
}

const (
    TenantContextName = "Tenant-Context"
    TenantHeader = "Tenant"
	TenantContextLevel  = -100
	)

func NewTenantContextObject(tenant string) TenantContextObject {
	return TenantContextObject{
		tenant: tenant,
	}
}

func (contextObject TenantContextObject) GetTenant() string {
	return contextObject.tenant
}

func (contextObject TenantContextObject) Serialize() (map[string]string, error) {
	if contextObject.tenant != "" {
		return map[string]string{TenantHeader: contextObject.tenant}, nil
	} else {
		return nil, nil
	}
}

func Of(ctx context.Context) (*TenantContextObject, error) {
	var contextObject TenantContextObject
	provider, err := ctxmanager.GetProvider(TenantContextName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s provider", TenantContextName)
	}
	abstractContextObject := provider.Get(ctx)
	if abstractContextObject == nil {
		return nil, errors.New("tenant context object is null")
	}
	contextObject = abstractContextObject.(TenantContextObject)
	return &contextObject, nil
}
