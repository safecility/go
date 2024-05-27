package stream

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// PublishToTopic just wraps json marshalling of an interface and Publish to a google pubsub Topic returning the result
// of the Publish method
func PublishToTopic(m interface{}, topic *pubsub.Topic) (*string, error) {
	if topic == nil {
		return nil, fmt.Errorf("no topic configured")
	}
	ctx := context.Background()

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not marshall message %v: %v", m, err)
	}

	message := &pubsub.Message{
		Data:        jsonBytes,
		PublishTime: time.Now(),
	}

	result, err := topic.Publish(ctx, message).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not publish json to topic %v: %v", topic, err)
	}

	return &result, nil
}

// GetDefaultSubscriptionConfig TODO fill in better defaults - make duration deployment dependant
func GetDefaultSubscriptionConfig(topic *pubsub.Topic, retentionDuration time.Duration) pubsub.SubscriptionConfig {

	return pubsub.SubscriptionConfig{
		Topic:                         topic,
		PushConfig:                    pubsub.PushConfig{},
		BigQueryConfig:                pubsub.BigQueryConfig{},
		CloudStorageConfig:            pubsub.CloudStorageConfig{},
		AckDeadline:                   0,
		RetainAckedMessages:           false,
		RetentionDuration:             retentionDuration,
		ExpirationPolicy:              nil,
		Labels:                        nil,
		EnableMessageOrdering:         false,
		DeadLetterPolicy:              nil,
		Filter:                        "",
		RetryPolicy:                   nil,
		Detached:                      false,
		TopicMessageRetentionDuration: 0,
		EnableExactlyOnceDelivery:     false,
		State:                         0,
	}
}
