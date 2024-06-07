package messages

import "github.com/safecility/go/lib/stream"

type MqttPath string

const (
	Join   MqttPath = "join"
	Uplink MqttPath = "uplink"
	Queued MqttPath = "down/queued"
	Ack    MqttPath = "down/ack"
	Nack   MqttPath = "down/nack"
	Sent   MqttPath = "down/sent"
	Failed MqttPath = "down/failed"

	AllDevices = "+"
	//Push   MqttPath = "down/push"
)

type Signal struct {
	Rssi int
	Snr  float64
}

type DeviceSignal struct {
	DeviceUID string
	Signal    Signal
}

type Location struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

type DeviceLocation struct {
	DeviceUID string
	Location  Location
}

type LoraMessage struct {
	stream.SimpleMessage
	Signal   *Signal
	Location *Location
}
