package client

import "testing"

func TestMgcUrl_String(t *testing.T) {
	tests := []struct {
		name string
		m    MgcUrl
		want string
	}{
		{
			name: "BrNe1 region string conversion",
			m:    BrNe1,
			want: "https://api.magalu.cloud/br-ne1",
		},
		{
			name: "BrSe1 region string conversion",
			m:    BrSe1,
			want: "https://api.magalu.cloud/br-se1",
		},
		{
			name: "BrMgl1 region string conversion",
			m:    BrMgl1,
			want: "https://api.magalu.cloud/br-se-1",
		},
		{
			name: "Empty URL string conversion",
			m:    "",
			want: "",
		},
		{
			name: "Custom URL string conversion",
			m:    "https://custom.url/",
			want: "https://custom.url/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("MgcUrl.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegionConstants(t *testing.T) {
	if BrNe1 != "https://api.magalu.cloud/br-ne1" {
		t.Errorf("BrNe1 constant has unexpected value: %s", BrNe1)
	}

	if BrSe1 != "https://api.magalu.cloud/br-se1" {
		t.Errorf("BrSe1 constant has unexpected value: %s", BrSe1)
	}
}
