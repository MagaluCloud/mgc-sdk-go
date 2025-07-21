# Helpers

FUNCTIONS
```
func BoolPtr(b bool) *bool
```
```
func Float32Ptr(f float32) *float32
```
```
func Float64Ptr(f float64) *float64
```
```
func Int16Ptr(i int16) *int16
```
```
func Int32Ptr(i int32) *int32
```
```
func Int64Ptr(i int64) *int64
```
```
func Int8Ptr(i int8) *int8
```
```
func IntPtr(i int) *int
```
```
func StrPtr(s string) *string
```
```
func Uint16Ptr(u uint16) *uint16
```
```
func Uint32Ptr(u uint32) *uint32
```
```
func Uint64Ptr(u uint64) *uint64
```
```
func Uint8Ptr(u uint8) *uint8
```
```
func UintPtr(u uint) *uint


```
```
type QueryParams interface {
Add(name string, value *string)
AddReflect(name string, value any)
Encode() string
}

```
```
func NewQueryParams(httpReq *http.Request) QueryParams


```

