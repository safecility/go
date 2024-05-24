package stream

import "time"

type BrokerDevice struct {
	Source    string
	DeviceUID string
}

type SimpleMessage struct {
	BrokerDevice
	Payload []byte
	Time    time.Time
}
