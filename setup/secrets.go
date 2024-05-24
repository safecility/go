package setup

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

// GetSecret retrieve a secret with the given version name
func GetSecret(versionName string) string {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup client")
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg("could not close client")
		}
	}(client)

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: versionName,
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatal().Err(err).Msg("could not access secret")
	}

	return fmt.Sprintf("%s", result.Payload.Data)
}
