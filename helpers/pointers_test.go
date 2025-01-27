package helpers

import "testing"

func TestIntPtr(t *testing.T) {
	value := 42
	ptr := IntPtr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestStrPtr(t *testing.T) {
	value := "test"
	ptr := StrPtr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestBoolPtr(t *testing.T) {
	value := true
	ptr := BoolPtr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestFloat32Ptr(t *testing.T) {
	value := float32(3.14)
	ptr := Float32Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestFloat64Ptr(t *testing.T) {
	value := 3.14159
	ptr := Float64Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestInt8Ptr(t *testing.T) {
	value := int8(8)
	ptr := Int8Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestInt16Ptr(t *testing.T) {
	value := int16(16)
	ptr := Int16Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestInt32Ptr(t *testing.T) {
	value := int32(32)
	ptr := Int32Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestInt64Ptr(t *testing.T) {
	value := int64(64)
	ptr := Int64Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestUintPtr(t *testing.T) {
	value := uint(42)
	ptr := UintPtr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestUint8Ptr(t *testing.T) {
	value := uint8(8)
	ptr := Uint8Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestUint16Ptr(t *testing.T) {
	value := uint16(16)
	ptr := Uint16Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestUint32Ptr(t *testing.T) {
	value := uint32(32)
	ptr := Uint32Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}

func TestUint64Ptr(t *testing.T) {
	value := uint64(64)
	ptr := Uint64Ptr(value)
	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}
	if *ptr != value {
		t.Errorf("Expected %v, got %v", value, *ptr)
	}
}
