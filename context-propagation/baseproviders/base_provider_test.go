package baseproviders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/configloader"
)

func init() {
	configloader.Init(configloader.EnvPropertySource())
}

func TestProviders(t *testing.T) {
	assert.Equal(t, 10, len(Get()))
}
