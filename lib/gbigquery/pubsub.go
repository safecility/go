package gbigquery

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

// CreateProtoSchema creates a schema resource from a schema proto file.
func CreateProtoSchema(client *pubsub.SchemaClient, schemaID, protoFile string) (*pubsub.SchemaConfig, error) {
	protoSource, err := os.ReadFile(protoFile)
	if err != nil {
		return nil, fmt.Errorf("error reading from file: %s", protoFile)
	}

	config := pubsub.SchemaConfig{
		Type:       pubsub.SchemaProtocolBuffer,
		Definition: string(protoSource),
	}

	ctx := context.Background()
	s, err := client.CreateSchema(ctx, schemaID, config)
	if err != nil {
		return nil, fmt.Errorf("CreateSchema: %w", err)
	}
	log.Debug().Str("schema", s.Name).Msg("Schema created")
	return s, nil
}

func CreateBigqueryTopic(client *pubsub.Client, topicName string, schema *pubsub.SchemaConfig) (*pubsub.Topic, error) {
	ctx := context.Background()
	bigqueryTopic, err := client.CreateTopicWithConfig(ctx, topicName, &pubsub.TopicConfig{
		SchemaSettings: &pubsub.SchemaSettings{
			Schema:   schema.Name,
			Encoding: pubsub.EncodingBinary,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("setup could not create topic")
	}
	log.Info().Str("topic", bigqueryTopic.String()).Msg("created topic")

	return bigqueryTopic, nil
}

// CreateBigQuerySubscription creates a Pub/Sub subscription that exports messages to BigQuery.
func CreateBigQuerySubscription(client *pubsub.Client, subscriptionName, table string, topic *pubsub.Topic) error {
	ctx := context.Background()

	sub, err := client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic: topic,
		BigQueryConfig: pubsub.BigQueryConfig{
			Table:             table,
			WriteMetadata:     false,
			UseTopicSchema:    true,
			DropUnknownFields: true,
		},
	})
	if err != nil {
		return err
	}
	log.Debug().Str("subscription", sub.ID()).Msg("Created BigQuery subscription")

	return nil
}
