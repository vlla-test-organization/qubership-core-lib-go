# How to create custom property source

Configloader library allows users to create their own custom property sources, so they can use different configuration 
sources, not only provided out of the box. Configloader allows defining as many extra property sources as user desires. 
This article provides an example of how it may be done.

- [Terms](#terms)
- [Process](#process)
- [Full example](#full-example)

## Terms

* Property source - any possible source of configuration parameters and properties, eg. configuration files, environment variables, parameters from config-server, etc.

## Process 

1. First step is to implement interfaces PropertyProvider and koanf.Parser from [koanf](https://github.com/knadh/koanf#readme).
   
   _PropertyProvider_ is used for read configuration from a source (file, HTTP, env, etc.)
    ```go
        type PropertyProvider interface {
            ReadBytes(k *koanf.Koanf) ([]byte, error)
            Read(k *koanf.Koanf) (map[string]interface{}, error)
        }
    ```
   * ReadBytes(*koanf.Koanf) returns the entire configuration as raw []bytes to be parsed with a koanf.Parser.
   * Read(*koanf.Koanf) returns the parsed configuration as a nested map[string]interface{}.
    
    During initialization koanf at first will try to use _Read(*koanf.Koanf)_ method and if it returns error, koanf will try to use ReadBytes(*koanf.Koanf).
    > Note #1 If user implemented Read(*koanf.Koanf) method there is no need to implement koanf.Parser interface, because Read(*koanf.Koanf) already provides map[string]interface{}. But if user implemented only ReadBytes(*koanf.Koanf) you have to implement koanf.Parser in order to parse bytes into map[string]interface{}.
    
    > Note #2 map[string]interface{} should be stored in flatten way in order to prevent incorrect properties overriding.

    > Note #3 If implementation of ReadBytes(*koanf.Koanf) or Read(*koanf.Koanf) methods need to get the value of a
    property at a given execution time, then it MUST use *koanf.Koanf instance from argument, NOT from exported
    functions like GetOrDefaultString etc. Invoking WRITE operations on that instance is forbidden.


    _koanf.Parser_ is used for parsing bytes, received from koanf.Provider.
     ```go
        type  Parser interface {
              Unmarshal([]byte) (map[string]interface{}, error)
              Marshal(map[string]interface{}) ([]byte, error)
        }
    ```
    * Unmarshal([]byte) takes bytes after ReadBytes() and parses them into map.

    > Note #3 Visit [documentation](https://github.com/knadh/koanf#bundled-providers) to find some embedded realisations for Parser and Provider.
2. Second step is to create instance of PropertySource struct
    ```go
        type PropertySource struct {
             Provider PropertyProvider
             Parser koanf.Parser
        }
    ```
   If user implemented only Provider with Read(*koanf.Koanf) method, they may pass nil instead of parser implementation.
3. Finally, user can add their custom property source to _configloader.Init_ method.
   
    Examples:
   
    Use implementation from koanf for both Parser and Provider.
    ```go
        func init() {
           newPropertySource := &PropertySource {
                Provider: AsPropertyProvider(file.Provider("testdata/test.json")),
                Parser:   json.Parser(),
           }
           configloader.Init([]*PropertySource{newPropertySource})
        }
    ```
   
    Use implementation from koanf, but without Parser
    ```go
        func init() {
           newPropertySource := &PropertySource {
                Provider: AsPropertyProvider(env.Provider("", ".", func(s string) string {return strings.Replace(strings.ToLower(s), "_", ".", -1)})),
                Parser:   nil,
           }
           configloader.Init([]*PropertySource{newPropertySource})
        }
    ```
   
## Full example

See [configserver-propertysource](https://github.com/netcracker/qubership-core-lib-go-rest-utils/blob/main/configserver-propertysource/README.md) as an example.
Here, configServerLoader implements _PropertyProvider_ interface. _Read(*koanf.Koanf)_ function uses argument as read source
and propagates this argument to all called functions / methods that needs to access already read properties, it's necessary because _*koanf.Koanf_ instance is not fully
loaded at execution time of Read function.


```go
type configServerLoader struct {
    propertySourceConfiguration *PropertySourceConfiguration
}
 
// Implement only Provider interface
func newConfigServerLoader(params *PropertySourceConfiguration) *configServerLoader {
    return &configServerLoader{params}
}
 
// Won't use raw bytes
func (this *configServerLoader) ReadBytes(_ *koanf.Koanf) ([]byte, error) {
    return nil, errors.New("configserver provider does not support this method")
}
 
// Will create map by ourselves
func (this configServerLoader) Read(k *koanf.Koanf) (map[string]interface{}, error) {
    source, err := getConfigServerProperties(k, this.propertySourceConfiguration)
    if err != nil {
        return nil, err
    }
    return maps.Unflatten(source, "."), nil
}
 
// some custom logic here
func getConfigServerProperties(konf *koanf.Koanf, params *PropertySourceConfiguration) (map[string]interface{}, error) {
    microserviceName, configServerUrl := getMicroserviceNameAndURL(konf, params)
    m2mToken, err := m2mmanager.GetToken(context.Background())
    client := &http.Client{}
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s/default", configServerUrl, microserviceName), nil)
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m2mToken))
    res, err := client.Do(req)
    if err != nil {
        log.Printf("Failed send requst to config-server: %s", err)
        return nil, err
    }
    defer res.Body.Close()
    return parseBody(res.Body)
}
 
func main() {
    configserverPropertySource := &configloader.PropertySource{Provider: newConfigServerLoader(&configuration)}
    // Init configloader with custom property source
    configloader.Init([]*PropertySource{configserverPropertySource })
}
```