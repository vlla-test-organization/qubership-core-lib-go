package configloader

import "github.com/knadh/koanf/v2"

// PropertyProvider purpose is the same as koanf.Provider, but it represents read behaviour that can be executed
// in dependent stages, so we can have access to already read elements from koanf.Koanf instance that is passed as argument.
// During properties initialization/refresh implementations can use Koanf instance in order to READ fresh properties from
// argument instead of invoking GetOrDefault and other functions. WRITE operations on Koanf instance is forbidden.
type PropertyProvider interface {
	ReadBytes(k *koanf.Koanf) ([]byte, error)
	Read(k *koanf.Koanf) (map[string]interface{}, error)
}

// koanfProviderAdapter adapts PropertyProvider interface to koanf.Provider interface
type koanfProviderAdapter struct {
	propertyProvider PropertyProvider
	konf             *koanf.Koanf
}

func (k *koanfProviderAdapter) ReadBytes() ([]byte, error) {
	return k.propertyProvider.ReadBytes(k.konf)
}

func (k *koanfProviderAdapter) Read() (map[string]interface{}, error) {
	return k.propertyProvider.Read(k.konf)
}

func asKoanfProvider(konf *koanf.Koanf, propProvider PropertyProvider) koanf.Provider {
	return &koanfProviderAdapter{
		propertyProvider: propProvider,
		konf:             konf,
	}
}

// PropertyProviderAdapter adapts koanf.Provider interface to PropertyProvider
type PropertyProviderAdapter struct {
	provider koanf.Provider
}

func (p *PropertyProviderAdapter) ReadBytes(_ *koanf.Koanf) ([]byte, error) {
	return p.provider.ReadBytes()
}

func (p *PropertyProviderAdapter) Read(_ *koanf.Koanf) (map[string]interface{}, error) {
	return p.provider.Read()
}

func AsPropertyProvider(provider koanf.Provider) PropertyProvider {
	return &PropertyProviderAdapter{provider: provider}
}
