package messages

import (
	"github.com/safecility/go/lib/stream"
)

type MqttChannel string

const (
	Join   MqttChannel = "join"
	Uplink MqttChannel = "uplink"
	Push   MqttChannel = "down/push"
	Queued MqttChannel = "down/queued"
	Ack    MqttChannel = "down/ack"
	Nack   MqttChannel = "down/nack"
	Sent   MqttChannel = "down/sent"
	Failed MqttChannel = "down/failed"

	AllDevices = "+"
)

type Signal struct {
	Rssi int
	Snr  float64
}

type LoraData struct {
	DeviceEUI []byte
	Signal    *Signal
	Channel   MqttChannel
}

// LoraMessage - enhance our basic message with Lora specific fields
type LoraMessage struct {
	LoraData
	stream.SimpleMessage
}
