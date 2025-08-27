package serviceloader

import (
	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v8/logging"
	"reflect"
	"sync"
	"testing"
)

type ServiceA interface {
	DoSomething() string
}

type ServiceB interface {
	DoSomethingElse() string
}

type ConcreteServiceA1 struct{}

func (s *ConcreteServiceA1) DoSomething() string {
	return "ConcreteServiceA1"
}

type ConcreteServiceA2 struct{}

func (s *ConcreteServiceA2) DoSomething() string {
	return "ConcreteServiceA2"
}

type ConcreteServiceB struct{}

func (s *ConcreteServiceB) DoSomethingElse() string {
	return "ConcreteServiceB"
}

func TestRegisterAndLoad(t *testing.T) {
	services = make([]serviceRegistration, 0)
	foundCache = &sync.Map{}

	Register(10, &ConcreteServiceA1{})
	Register(20, &ConcreteServiceA2{})
	Register(10, &ConcreteServiceB{})

	loadedServiceA, found := Load[ServiceA]()
	assert.True(t, found)
	assert.Equal(t, "ConcreteServiceA2", loadedServiceA.DoSomething())

	loadedServiceB, found := Load[ServiceB]()
	assert.True(t, found)
	assert.Equal(t, "ConcreteServiceB", loadedServiceB.DoSomethingElse())

	type NonExistentService interface {
		DoNothing() string
	}
	_, found = Load[NonExistentService]()
	assert.False(t, found)
}

func TestPriority(t *testing.T) {
	services = make([]serviceRegistration, 0)
	foundCache = &sync.Map{}

	Register(2, &ConcreteServiceA1{})
	Register(1, &ConcreteServiceA2{})

	loadedServiceA, found := Load[ServiceA]()
	assert.True(t, found)
	assert.Equal(t, "ConcreteServiceA1", loadedServiceA.DoSomething())
}

func TestConcurrentAccess(t *testing.T) {
	currentLvl := logger.GetLevel()
	logger.SetLevel(logging.LvlWarn)
	defer logger.SetLevel(currentLvl)

	services = make([]serviceRegistration, 0)
	foundCache = &sync.Map{}

	var wg sync.WaitGroup
	numGoroutines := 1000
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			Register(1, &ConcreteServiceA1{})
			Register(2, &ConcreteServiceA2{})
			loadedServiceA, found := Load[ServiceA]()
			assert.True(t, found)
			assert.Equal(t, "ConcreteServiceA2", loadedServiceA.DoSomething())
		}()
	}

	wg.Wait()
}

func TestMustLoad(t *testing.T) {
	services = make([]serviceRegistration, 0)
	foundCache = &sync.Map{}

	Register(10, &ConcreteServiceA1{})

	loadedServiceA := MustLoad[ServiceA]()
	assert.Equal(t, "ConcreteServiceA1", loadedServiceA.DoSomething())

	type NonExistentService interface {
		DoNothing() string
	}

	assert.Panicsf(t, func() { MustLoad[NonExistentService]() },
		"can not find implementation for '%s' in service loading context; make sure it was register before", reflect.TypeFor[NonExistentService](),
	)
}
