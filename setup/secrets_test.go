package setup

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
	"testing"
)

func cleanup(secretName string) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg("could not close secret manager client")
		}
	}(client)

	err = client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: secretName,
	})
	if err != nil {
		log.Err(err).Msg("could not delete secret")
	}
}

func TestGetNewSecrets(t *testing.T) {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg("could not close secret manager client")
		}
	}(client)

	type args struct {
		projectID string
		client    *secretmanager.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "TestGetNewSecrets",
			args:    args{projectID: "safecility-test", client: client},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNewSecrets(tt.args.projectID, tt.args.client); got == nil {
				t.Errorf("GetNewSecrets() = %v, want !nil", got)
			}
		})
	}
}

func TestSecrets_SetSecret(t *testing.T) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg("could not close secret manager client")
		}
	}(client)

	s := GetNewSecrets("safecility-test", client)

	type fields struct {
		secrets *Secrets
	}
	type args struct {
		secretID string
		payload  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *secretmanagerpb.Secret
		wantErr bool
	}{
		// TODO: Add more better test cases.
		{
			name: "TestSetSecret",
			fields: fields{
				secrets: s,
			},
			args: args{
				secretID: "test-new-secret",
				payload:  nil,
			},
			want:    &secretmanagerpb.Secret{Name: "test-new-secret"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.secrets.SetSecret(tt.args.secretID, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !strings.HasSuffix(got.GetName(), tt.want.Name) {
				t.Errorf("SetSecret() got = %v, want %v", got.GetName(), tt.want)
			}
			cleanup(got.Name)
			return
		})
	}
}

func TestSecrets_AddSecretVersion(t *testing.T) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg("could not close secret manager client")
		}
	}(client)

	type fields struct {
		projectID string
		client    *secretmanager.Client
	}
	type args struct {
		secretName string
		newPayload []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *secretmanagerpb.SecretVersion
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestAddSecretVersion",
			fields: fields{
				projectID: "safecility-test",
				client:    client,
			},
			args: args{
				secretName: "testsecret",
				newPayload: []byte("testpayload"),
			},
			want: &secretmanagerpb.SecretVersion{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secrets{
				client: tt.fields.client,
			}
			got, err := s.AddSecretVersion(tt.args.secretName, tt.args.newPayload)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSecretVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddSecretVersion() got = %v, want %v", got, tt.want)
			}
			cleanup(got.Name)
		})
	}
}

func TestSecrets_GetSecret(t *testing.T) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg("could not close secret manager client")
		}
	}(client)

	type fields struct {
		projectID string
		client    *secretmanager.Client
	}
	type args struct {
		secretName string
		version    uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestGetSecret",
			fields: fields{
				projectID: "safecility-test",
				client:    client,
			},
			args: args{
				secretName: "testsecret",
				version:    1,
			},
			want:    "testsecret",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secrets{
				client: tt.fields.client,
			}
			got, err := s.GetSecret(tt.args.secretName, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecrets_UpdateSecret(t *testing.T) {
	type fields struct {
		client *secretmanager.Client
	}
	type args struct {
		secretName string
		version    uint16
		values     map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Secrets{
				client: tt.fields.client,
			}
			if err := s.UpdateSecret(tt.args.secretName, tt.args.version, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
