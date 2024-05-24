package lib

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/mqtt/messages"
)

// UidTransformer allows paho UIDs to be transformed for different microservices/storage/cache implementations
type UidTransformer interface {
	GetUID(deviceID string) string
}

// PahoTransformer creates a SimpleMessage from a pahoMessage handling different LoRA data brokers
type PahoTransformer interface {
	UplinkTopic(deviceID string) string
	UplinkErrorsTopic(deviceID string) string
	MqttTopic(deviceID string, channel messages.MqttChannel) string
	GetDownlinkTopicPush(devID string) string
	GetDownlinkTopicReplace(devID string) string
	TransformPahoJoinMessage(m paho.Message) (*messages.LoraMessage, error)
	TransformPahoUplinkMessage(m paho.Message) (*messages.LoraMessage, error)
	TransformPahoDownlinkMessage(m paho.Message, channel messages.MqttChannel) (*messages.LoraMessage, error)
	CreateDownlink(message stream.SimpleMessage, correlationIDs []string) ([]byte, error)
}

// PayloadAdjuster differs from PahoTransformer by changing a SimpleMessage payload (to be Segments etc.)
// This allows us to deliver similar messages from LoRA to those from NBIoT allowing the same microservices to parse both
type PayloadAdjuster interface {
	AdjustPayload(m *messages.LoraMessage) error
}

type IdentityAdjuster struct{}

func (i IdentityAdjuster) AdjustPayload(_ *messages.LoraMessage) error {
	return nil
}
