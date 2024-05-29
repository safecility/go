package setup

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"reflect"
	"testing"
)

func TestMySQLConfig_connectionString(t *testing.T) {
	type fields struct {
		Password               string
		Username               string
		Host                   string
		Port                   int
		Database               string
		InstanceConnectionName string
	}
	type args struct {
		databaseName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MySQLConfig{
				Password:               tt.fields.Password,
				Username:               tt.fields.Username,
				Host:                   tt.fields.Host,
				Port:                   tt.fields.Port,
				Database:               tt.fields.Database,
				InstanceConnectionName: tt.fields.InstanceConnectionName,
			}
			if got := c.connectionString(tt.args.databaseName); got != tt.want {
				t.Errorf("connectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSafecilitySql(t *testing.T) {
	type args struct {
		config MySQLConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSafecilitySql(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSafecilitySql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSafecilitySql() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisConfig_Address(t *testing.T) {
	type fields struct {
		Host string
		Port string
		Key  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisConfig{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
				Key:  tt.fields.Key,
			}
			if got := r.Address(); got != tt.want {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisConfig_NewClient(t *testing.T) {
	type fields struct {
		Host string
		Port string
		Key  string
	}
	tests := []struct {
		name   string
		fields fields
		want   *redis.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisConfig{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
				Key:  tt.fields.Key,
			}
			if got := r.NewClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
