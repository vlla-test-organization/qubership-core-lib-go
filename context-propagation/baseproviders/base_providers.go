package baseproviders

import (
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/acceptlanguage"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/allowedheaders"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/apiversion"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/businessprocess"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/clientip"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/originatingbiid"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/tenant"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/xrequestid"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/xversion"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/baseproviders/xversionname"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/context-propagation/ctxmanager"
)

func Get() []ctxmanager.ContextProvider {
	return []ctxmanager.ContextProvider{
		acceptlanguage.AcceptLanguageProvider{},
		xversion.XVersionProvider{},
		xversionname.XVersionNameProvider{},
		apiversion.ApiVersionProvider{},
		xrequestid.XRequestIdProvider{},
		allowedheaders.NewAllowedHeaderProvider(),
		businessprocess.BusinessProcessProvider{},
		originatingbiid.OriginatingBiIdProvider{},
		clientip.ClientIpProvider{},
		tenant.TenantProvider{},
	}
}
