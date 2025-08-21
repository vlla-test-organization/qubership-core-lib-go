package serviceloader

import (
	"fmt"
	"github.com/vlla-test-organization/qubership-core-lib-go/v6/logging"
	"reflect"
	"sync"
)

var logger = logging.GetLogger("serviceloader")

type serviceRegistration struct {
	priority int
	instance any
}

var services = make([]serviceRegistration, 0)
var foundCache = &sync.Map{}
var mtx = sync.RWMutex{}

func Register[T any](priority int, service *T) {
	mtx.Lock()
	defer mtx.Unlock()
	foundCache.Clear() // reset cache
	services = append(services, serviceRegistration{priority: priority, instance: service})
}

func MustLoad[T any]() T {
	instance, found := Load[T]()
	if !found {
		panic(fmt.Sprintf("can not find implementation for `%s' in service loading; make sure it was register before", reflect.TypeFor[T]()))
	}
	return instance
}

func Load[T any]() (T, bool) {
	targetType := reflect.TypeFor[T]()
	if instance, ok := foundCache.Load(targetType); ok {
		return instance.(T), true
	}

	mtx.RLock()
	defer mtx.RUnlock()

	found := serviceRegistration{}
	for _, reg := range services {
		if reflect.TypeOf(reg.instance).Implements(targetType) && (found.instance == nil || found.priority < reg.priority) {
			found = reg
		}
	}

	if found.instance != nil {
		logger.Info("Located `%s' as implementation for `%s'", reflect.TypeOf(found.instance), reflect.TypeFor[T]())
		foundCache.Store(targetType, found.instance)
		return found.instance.(T), true
	} else {
		return *new(T), false
	}
}
