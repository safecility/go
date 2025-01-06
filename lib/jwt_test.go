package lib

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
	"testing"
	"time"
)

func TestJWTParser_CreateToken(t *testing.T) {
	type fields struct {
		hmacSecret []byte
	}
	type args struct {
		claims jwt.MapClaims
	}
	now := time.Now()
	expires := now.Add(1 * time.Hour)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *jwt.MapClaims
		wantErr bool
	}{
		// TODO: Add more test cases.
		{
			name: "basic claims",
			fields: fields{
				hmacSecret: []byte("a secret"),
			},
			args: args{
				claims: jwt.MapClaims{
					"name":    "claimWithExpire",
					"value":   "claim value 1",
					"created": now.Format(time.RFC3339),
					"exp":     expires.Unix(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JWTParser{
				hmacSecret: tt.fields.hmacSecret,
			}
			got, err := p.CreateToken(tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			token, err := p.ParseToken(got)
			if err != nil {
				t.Errorf("CreateToken() must be reversible, error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(token)
		})
	}
}

func TestJWTParser_ParseToken(t *testing.T) {
	type fields struct {
		hmacSecret []byte
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *jwt.MapClaims
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JWTParser{
				hmacSecret: tt.fields.hmacSecret,
			}
			got, err := p.ParseToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJWTParser(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
		want JWTParser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJWTParser(tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJWTParser() = %v, want %v", got, tt.want)
			}
		})
	}
}
