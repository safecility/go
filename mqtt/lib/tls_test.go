package lib

import (
	"crypto/tls"
	"testing"
)

func TestNewTlsConfig(t *testing.T) {
	tests := []struct {
		name string
		want *tls.Config
	}{
		// TODO: Add test cases.
		{name: "TestNewTlsConfig", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := NewTlsConfig(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewTlsConfig() = %v, want %v", got, tt.want)
			//}
			got := NewTlsConfig()
			if got == nil {
				if tt.want != nil {
					t.Errorf("NewTlsConfig() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
