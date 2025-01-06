package setup

import (
	"fmt"
	"testing"
)

type TestConfig struct {
	Test struct {
		Name   string   `json:"name"`
		Array  []string `json:"array"`
		Object struct {
			Name     string
			Position int
			Value    float32
		}
	}
}

func TestGetConfig(t *testing.T) {
	type args struct {
		deployment string
		config     any
	}
	existingConfig := &TestConfig{}
	existingTestFile := fmt.Sprintf("./testdata/%s", Test)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add more test cases.
		{
			name: "exists",
			args: args{
				deployment: existingTestFile,
				config:     existingConfig,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetConfig(tt.args.deployment, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.args.config)
		})
	}
}
