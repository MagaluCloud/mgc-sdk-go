package helpers

//go:fix inline
func IntPtr(i int) *int {
	return new(i)
}

//go:fix inline
func StrPtr(s string) *string {
	return new(s)
}

//go:fix inline
func BoolPtr(b bool) *bool {
	return new(b)
}

//go:fix inline
func Float32Ptr(f float32) *float32 {
	return new(f)
}

//go:fix inline
func Float64Ptr(f float64) *float64 {
	return new(f)
}

//go:fix inline
func Int8Ptr(i int8) *int8 {
	return new(i)
}

//go:fix inline
func Int16Ptr(i int16) *int16 {
	return new(i)
}

//go:fix inline
func Int32Ptr(i int32) *int32 {
	return new(i)
}

//go:fix inline
func Int64Ptr(i int64) *int64 {
	return new(i)
}

//go:fix inline
func UintPtr(u uint) *uint {
	return new(u)
}

//go:fix inline
func Uint8Ptr(u uint8) *uint8 {
	return new(u)
}

//go:fix inline
func Uint16Ptr(u uint16) *uint16 {
	return new(u)
}

//go:fix inline
func Uint32Ptr(u uint32) *uint32 {
	return new(u)
}

//go:fix inline
func Uint64Ptr(u uint64) *uint64 {
	return new(u)
}
