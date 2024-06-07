package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/mqtt/messages"
	"time"
)

// TtnV3 simplifies access to Uplink and Downlink pubsubs
// UidTransformer can be nil in which case the DeviceID is preserved - however, an AppIDTransformer is usually preferred
type TtnV3 struct {
	AppID          string
	UidTransformer UidTransformer
}

// UplinkTopic can be replaced by calling MqttTopic with channel Up
func (t TtnV3) UplinkTopic(deviceID string) string {
	return fmt.Sprintf(`v3/%s@%s/devices/%s/up`, t.AppID, "ttn", deviceID)
}

// UplinkErrorsTopic -listen for device errors
func (t TtnV3) UplinkErrorsTopic(deviceID string) string {
	return fmt.Sprintf(`v3/%s@%s/devices/%s/events/up/error`, t.AppID, "ttn", deviceID)
}

// GetDownlinkTopicPush can be replaced by calling MqttTopic with channel Push
func (t TtnV3) GetDownlinkTopicPush(devID string) string {
	return fmt.Sprintf("v3/%s@%s/devices/%s/down/push", t.AppID, "ttn", devID)
}
func (t TtnV3) GetDownlinkTopicReplace(devID string) string {
	return fmt.Sprintf("v3/%s@%s/devices/%s/down/replace", t.AppID, "ttn", devID)
}
func (t TtnV3) MqttTopic(deviceID string, channel messages.MqttPath) string {
	return fmt.Sprintf(`v3/%s@%s/devices/%s/%s`, t.AppID, "ttn", deviceID, channel)
}

func (t TtnV3) TransformPahoUplinkMessage(m paho.Message) (*messages.LoraMessage, error) {
	uplink := &UplinkV3{}
	err := json.Unmarshal(m.Payload(), uplink)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("decode", fmt.Sprintf("%s", m.Payload())).
		Str("uplink", fmt.Sprintf("%+v", uplink)).Msg("payload")

	deviceID := uplink.MessageIDs.DeviceID
	if t.UidTransformer != nil {
		deviceID = t.UidTransformer.GetUID(deviceID)
	}
	payload, err := base64.StdEncoding.DecodeString(uplink.Message.Payload)
	if err != nil {
		return nil, err
	}

	bd := stream.BrokerDevice{
		Source:    t.AppID,
		DeviceUID: deviceID,
	}

	sm := stream.SimpleMessage{
		BrokerDevice: bd,
		Payload:      payload,
		Time:         uplink.Received,
	}

	lm := &messages.LoraMessage{
		SimpleMessage: sm,
	}

	if len(uplink.Message.Metadata) > 0 {
		md := uplink.Message.Metadata[0]

		lm.Location = &messages.Location{
			Latitude:  md.Location.Latitude,
			Longitude: md.Location.Longitude,
			Altitude:  float64(md.Location.Altitude),
		}
		lm.Signal = &messages.Signal{
			Rssi: md.Rssi,
			Snr:  md.Snr,
		}
	}

	return lm, nil
}

func (t TtnV3) TransformPahoJoinMessage(m paho.Message) (*stream.SimpleMessage, error) {
	join := &JoinV3{}
	err := json.Unmarshal(m.Payload(), join)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("payload", fmt.Sprintf("decoded %s", m.Payload())).Msg("join")
	tm := join.Received
	if tm.IsZero() {
		tm = time.Now()
	}
	deviceID := join.MessageIDs.DeviceID
	if t.UidTransformer != nil {
		deviceID = t.UidTransformer.GetUID(deviceID)
	}
	payload, err := base64.StdEncoding.DecodeString(join.Payload)
	if err != nil {
		return nil, err
	}

	bd := stream.BrokerDevice{
		Source:    t.AppID,
		DeviceUID: deviceID,
	}

	sm := &stream.SimpleMessage{
		BrokerDevice: bd,
		Payload:      payload,
		Time:         tm,
	}

	return sm, nil
}

func (t TtnV3) TransformPahoDownlinkMessage(m paho.Message, path messages.MqttPath) (*stream.SimpleMessage, error) {
	downlink := &DownlinkV3{}
	err := json.Unmarshal(m.Payload(), downlink)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("payload", fmt.Sprintf("%s", m.Payload())).Msg("downlink")
	var p = ""

	switch path {
	case messages.Queued:
		p = downlink.DownlinkQueued.Payload
		break
	case messages.Sent:
		p = downlink.DownlinkSent.Payload
		break
	case messages.Failed:
		p = downlink.DownlinkFailed.Payload
		break
	}

	payload, err := base64.StdEncoding.DecodeString(p)
	if err != nil {
		log.Log().Msg("payload is invalid")
		return nil, err
	}

	tm := time.Now()

	bd := stream.BrokerDevice{
		Source:    t.AppID,
		DeviceUID: downlink.MessageIDs.DeviceID,
	}

	sm := &stream.SimpleMessage{
		BrokerDevice: bd,
		Payload:      payload,
		Time:         tm,
	}

	return sm, nil
}

func (t TtnV3) CreateDownlink(message stream.SimpleMessage, correlationIDs []string) ([]byte, error) {
	dl := t.createV3Downlink(message, correlationIDs)
	return json.Marshal(dl)
}

func (t TtnV3) createV3Downlink(message stream.SimpleMessage, correlationIDs []string) *DownlinksV3 {
	encoded := base64.StdEncoding.EncodeToString(message.Payload)

	return &DownlinksV3{Downlinks: []DownlinkV3Body{{
		FPort:          15,
		Payload:        encoded,
		Priority:       "HIGH",
		Confirmed:      true,
		CorrelationIDs: correlationIDs,
	}}}
}

type DownlinksV3 struct {
	Downlinks []DownlinkV3Body `json:"downlinks"`
}

type DownlinkV3Body struct {
	FPort          int      `json:"f_port"`
	Payload        string   `json:"frm_payload"`
	Priority       string   `json:"priority"`
	Confirmed      bool     `json:"confirmed"`
	CorrelationIDs []string `json:"correlation_ids"`
}

type JoinV3 struct {
	AppID      string
	MessageIDs MessageIDs `json:"end_device_ids"`
	Received   time.Time  `json:"received_at"`
	Payload    string     `json:"frm_payload"`
}

type DownlinkV3 struct {
	AppID          string
	MessageIDs     MessageIDs     `json:"end_device_ids"`
	DownlinkQueued DownlinkV3Body `json:"downlink_queued"`
	DownlinkSent   DownlinkV3Body `json:"downlink_sent"`
	DownlinkAck    DownlinkV3Body `json:"downlink_ack"`
	DownlinkFailed DownlinkV3Body `json:"downlink_failed"`
}

type UplinkV3 struct {
	AppID      string
	Message    UplinkMessage `json:"uplink_message"`
	MessageIDs MessageIDs    `json:"end_device_ids"`
	Received   time.Time     `json:"received_at"`
}

type RxMetadata struct {
	GatewayIds struct {
		GatewayId string `json:"gateway_id"`
		Eui       string `json:"eui"`
	} `json:"gateway_ids"`
	Time        time.Time `json:"time"`
	Timestamp   int       `json:"timestamp"`
	Rssi        int       `json:"rssi"`
	ChannelRssi int       `json:"channel_rssi"`
	Snr         float64   `json:"snr"`
	Location    struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Altitude  int     `json:"altitude"`
		Source    string  `json:"source"`
	} `json:"location"`
	UplinkToken string `json:"uplink_token"`
}

type UplinkMessage struct {
	Payload        string       `json:"frm_payload"`
	DecodedPayload interface{}  `json:"decoded_payload"`
	Metadata       []RxMetadata `json:"rx_metadata"`
}

type MessageIDs struct {
	DeviceID  string `json:"device_id"`
	DeviceEUI string `json:"dev_eui"`
}
