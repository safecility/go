package setup

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/protobuf/field_mask"
)

type Secret struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type Secrets struct {
	projectID string
	client    *secretmanager.Client
}

// GetNewSecrets - fix all the weird things in google's api
func GetNewSecrets(projectID string, client *secretmanager.Client) *Secrets {
	return &Secrets{projectID: projectID, client: client}
}

func (s *Secrets) Close() error {
	return s.client.Close()
}

func (s *Secrets) SetSecret(secretID string, payload []byte) (*secretmanagerpb.Secret, error) {
	ctx := context.Background()

	// Create the request to create the secret.
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", s.projectID),
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	secret, err := s.client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		return nil, err
	}

	if payload != nil {
		var version *secretmanagerpb.SecretVersion
		version, err = s.AddSecretVersion(secretID, payload)
		if version != nil {
			log.Debug().Str("version", version.Name)
		}
	}

	return secret, err
}

func (s *Secrets) AddSecretVersion(secretName string, newPayload []byte) (*secretmanagerpb.SecretVersion, error) {
	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("projects/%s/secrets/%s", s.projectID, secretName),
		Payload: &secretmanagerpb.SecretPayload{
			Data: newPayload,
		},
	}

	ctx := context.Background()
	// Call the API.
	return s.client.AddSecretVersion(ctx, req)
}

// GetSecret retrieve a secret with the given version name.
// The version name must comply with naming convention - auto is 1, 2 etc
func (s *Secrets) GetSecret(secret Secret) ([]byte, error) {
	// Create the client.
	ctx := context.Background()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%d", s.projectID, secret.Name, secret.Version),
	}

	result, err := s.client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, err
	}

	return result.Payload.Data, nil
}

func (s *Secrets) UpdateSecret(secret Secret, labels map[string]string) error {
	// Build the request.
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%d", s.projectID, secret.Name, secret.Version)
	req := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name:   name,
			Labels: labels,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"labels"},
		},
	}

	ctx := context.Background()
	// Call the API.
	result, err := s.client.UpdateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	log.Debug().Str("secret", result.Name).Msg("updated secret")
	return nil
}
