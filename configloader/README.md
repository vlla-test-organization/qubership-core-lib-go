# Configloader-base

This package allows loading configuration properties by provided property sources. Out of the package we provide two basic 
provider `yaml property source` and `environment property source` but we have a possibility to provide your own. 
This package is based on the [koanf](https://github.com/knadh/koanf) library.

- [Installation](#installation)
- [Important notes](#important-notes)
- [Usage](#usage)
- [Available property sources](#available-property-sources)
- [How to create custom property source](#how-to-create-custom-property-source)  
- [Quick example](#quick-example)

## Installation

To install `configloader` use
```go
 go get github.com/netcracker/qubership-core-lib-go/v3@<latest released version>
```

## Usage

First user have to initialize the configloader. Configloader provides two methods for initialization: 
* `configloader.InitWithSourcesArray(propertySources []*configloader.PropertySource)` 
* `configloader.Init(propertySources ...*PropertySource)`

Available sources are defined in [Available Property Sources](#available-property-sources).
```go
configloader.Init(propSource1, propSource2, propSource3)
// or
propSourcesArray := []*configloader.PropertySource{propSource1, propsSource2, propSource3}
configloader.InitWithSourcesArray(propSourcesArray)
// or with built-in func
configloader.InitWithSourcesArray(configloader.BasePropertySource())
```
After that user can access values by decorated methods: `GetOrDefault` or `GetOrDefaultString` or can use `GetKoanf()` for direct dealing with Koanf object.

`GetOrDefault(key string, def interface{})` allows to get any object (for example, map or slice) from property source. 
* key - is a name of property with dots (".") as separator (example: log.level.package). 
* def - default value, which should be returned if value with such key is missing.

`GetOrDefaultString(key string, def string)` allows to get a string value from property source.
* key - is a name of property with dots (".") as separator (example: log.level or LOG_LEVEL).
* def - default value, which should be returned if value with such key is missing.

_Please note, config keys are case-sensitive in configloader. For example, `app.server.port` and `APP.SERVER.port` are not the same._
_To get environment variables use lowercase with dots as separator._

**Important** Property source which has a higher order in a specified list in `configloader.Init([]*configloader.PropertySource{})`
method takes precedence over property sources which are to the left (have less order number)
  ```go
   func TestCreateCustomPropertySource(t *testing.T) {
        b := []byte(`{"key": "bytes"}`)
        bytePropertySource := &PropertySource{
          Provider: AsPropertyProvider(rawbytes.Provider(b)),
          Parser:   json.Parser(),
        }
        
        jsonPropertySource := &PropertySource{
          Provider: AsPropertyProvider(file.Provider("testdata/test.json")), // {"key":"json"}
          Parser:   json.Parser(),
        }
        
        // bytePropertySource has higher priority than jsonPropertySource
        Init(jsonPropertySource, bytePropertySource)
        assert.Equal(t, "bytes", GetKoanf().Get("key"))
        assert.NotEqual(t, "json", GetKoanf().Get("key"))
    }
  ```

If content from any property source is changed, then you can re-read configuration properties from there. For such scenario,
use ```Refresh()``` function.
* Refresh will use ```property sources``` that is passed previously to ```Init``` function.
* ```Refresh``` will return error if it was called before ```Init``` function.
```go
  func TestRefreshRewritesPropertiesOnEnvVariables(t *testing.T) {
      os.Setenv("TEST_REFRESH_REWRITES_PROPERTIES", "old-value")
      
      InitWithSourcesArray([]*PropertySource{EnvPropertySource()})
      gotValue := GetOrDefault("test.refresh.rewrites.properties", nil)
      assert.Equal(t, "old-value", gotValue)
      
      os.Setenv("TEST_REFRESH_REWRITES_PROPERTIES", "new-value") // Property source content changed (env value changed)
      
      assert.Nil(t, Refresh()) // Re-read properties from the same EnvPropertySource
      gotValue = GetOrDefault("test.refresh.rewrites.properties", nil)
      assert.Equal(t, "new-value", gotValue)
  }

```

## Available Property Sources

There is a group of predefined functions which user can use as a parameter for `configloader.Init(propertySources []*PropertySource)`.

* Environment - use environment variables. 
  ```go
    configloader.EnvPropertySource()
  ```
* Yaml - use variables defined in `application.yaml`. 
  
  There are several options to specify the location of yaml file.
  * With `YamlPropertySourceParams` struct by passing whole path to new file with any name (eg. "_../data/properties.yaml_").
  * With environment variable `PROPERTY_FILE_PATH` by passing only path to folder where `application.yaml` file is located.
  * To use default path just put `application.yaml` file in root folder near _main_.
 
  ```go
    // if application.yaml located in "../resources/data/application.yaml"
    configloader.YamlPropertySource(configloader.YamlPropertySourceParams{ConfigFilePath: "../resources/data/application.yaml"})
    
    // or if application.yaml located in root
    configloader.YamlPropertySource()
  
    // or by environment, if application.yaml located in "/app/application.yaml"
    os.Setenv(PROPERTY_FILE_PATH, "/app/")
    configloader.YamlPropertySource()
  ```

* BasePropertySources - it's a convenient way to get both environment and yaml property source in the correct order
  ( environment property source has higher priority than yaml property source.):
  ```go
    configloader.BasePropertySources()
       
    // or if you want to specify path
     configloader.BasePropertySources(configloader.YamlPropertySourceParams{ConfigFilePath: "../resources/data/application.yaml"})
  ```

## How to create custom property source

Here you can find brief description. If you need more information see [this article](./docs/custom-property-source.md)

User can define its own custom Property Source. To do this, just create an instance of the `PropertySource` structure and use it as a parameter for `configloader.Init()`. 
Struct `PropertySource` has fields: `Provider` of type `PropertyProvider` and `Parser` of type `koanf.Parser`.

User have to implement `koanf.Parser` and `PropertyProvider` interfaces or just use ready implementations from koanf
(there are ready implementations for vault, json, rawbytes, etc) and wrap them using PropertyProviderAdapter. 

```go
    newPropertySource := &PropertySource {
        Provider: some implemented PropertyProvider,
        Parser:   some implemented koanf.Parser,
    }
```
You can use ready koanf.Provider implementations in following way:
```go
    newPropertySource := &PropertySource {
        Provider: AsPropertyProvider(/*pass here some implementation of koanf.Provider interface*/),
        Parser:   some implemented koanf.Parser,
    }
```

> Please note, that functions `Read` and `ReadBytes` from PropertyProvider `must return flatten maps with values`.
> All PropertySources should store properties in the same way in order to avoid undetermined behavior.

Example:
```go
    newPropertySource := &PropertySource {
        Provider: AsPropertyProvider(file.Provider("testdata/test.json")),
        Parser:   json.Parser(),
    }
    configloader.Init(newPropertySource)
```

## Subscription on config changes
Client code can un/subscribe on configuration changes via `Subscribe(EventHandler) (SubscriptionId, error)` and
`Unsubscribe(SubscriptionId) error` functions. For each `Refresh()` and `Init()` invocation all event handlers will be 
executed in async way. Event handler is a function with following declaration:
```go
type EventHandler func(Event) error
    
type Event struct {
    Type EventT
    Data interface{}
}
```
Each event contains type of event and supplied data. Exact data type corresponds to exact event type so client code must
perform type cast (and check for `nil` value). If client code want to unregister event handler, then it must invoke
`Unsubscribe` with parameter value obtained previously from `Subscribe` function result.

Here is a dumb usage example:
```go
func ExampleSubscribe() {
    handlerF := func(e Event) error {
        if e.Type == InitedEventT {
            fmt.Println("Config inited")
        } else if e.Type == RefreshedEventT {
            fmt.Println("Config refreshed")
        } else {
            fmt.Println("Unknown event sent")
        }
        return nil
    }
    id, err := Subscribe(handlerF)
    if err != nil {
        // error handling
    }

    _ = Refresh()
    // Output: Config refreshed

    if err := Unsubscribe(id); err != nil {
        // error handling
    }
}
```

## Quick example

application.yaml 
```yaml
application:
  name: example
slice:
  - first
  - second
```

Environment variables
```
  ENV_VAR=env
```

application

```go
package main

import (
  "fmt"
  "github.com/netcracker/qubership-core-lib-go/v3/configloader"
)

func init() {
    configloader.InitWithSourcesArray(configloader.BasePropertySources())
}

func main() {
  name := configloader.GetOrDefaultString("application.name", "empty")
  defaultValue := configloader.GetOrDefaultString("application.empty", "empty")
  sliceVar := configloader.GetOrDefault("slice", []string{"test"})
  envVar := configloader.GetOrDefaultString("env.var", "")
  fmt.Println(name)
  fmt.Println(defaultValue)
  fmt.Println(sliceVar)
  fmt.Println(envVar)
}
```

Output will be:
```
    example
    empty
    [first second]
    env
```
