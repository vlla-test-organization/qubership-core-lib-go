## Deprecation of configloader.GetEventChannel function

GetEventChannel return type do not allow to provide event bus behaviour, because channels are meant to be used in 2
ways:

- When there is N senders and 1 consumer
- When we close channel in order to notify multiple consumers. This solution lacks of possibility of
  sending event payload

### Solution

Introduce new callback-based (Observer pattern) functionality in order to react on config initializing / refreshing.
Client code must use `Subscribe` function to register callback. Here is example code:

```go
func ExampleSubscribe() {
handlerF := func (e Event) error {
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

### What should be done before major release

- Remove `configloader.INITED_EVENT` and `configloader.REFRESHED_EVENT`.
- Remove `eventChannel` in configloader.go.
