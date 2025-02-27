package constants

import (
	tlsUtils "github.com/netcracker/qubership-core-lib-go/v3/utils"
)

const (
	ConfigServerUrlProperty    = "config-server.url"
	httpDefaultConfigServerUrl = "http://config-server:8080"

	NamespaceProperty = "microservice.namespace"
	DefaultNamespace  = "unknown"

	MicroserviceNameProperty = "microservice.name"
	DefaultMicroserviceName  = "unknown"

	ProfileProperty = "profile"
	DefaultProfile  = "default"

	ServerHostnameProperty = "cloud.public.host"
	DefaultServerHostname  = "unknown"
)

func SelectUrl(httpUrl string, httpsUrl string) string {
	if tlsUtils.IsTlsEnabled() {
		return httpsUrl
	} else {
		return httpUrl
	}
}
