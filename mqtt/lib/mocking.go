package lib

import "sync"

type mockMessage struct {
	duplicate bool
	qos       byte
	retained  bool
	topic     string
	messageID uint16
	payload   []byte
	once      sync.Once
	ack       func()
}

func (m *mockMessage) Duplicate() bool {
	return m.duplicate
}

func (m *mockMessage) Qos() byte {
	return m.qos
}

func (m *mockMessage) Retained() bool {
	return m.retained
}

func (m *mockMessage) Topic() string {
	return m.topic
}

func (m *mockMessage) MessageID() uint16 {
	return m.messageID
}

func (m *mockMessage) Payload() []byte {
	return m.payload
}

func (m *mockMessage) Ack() {
	m.once.Do(m.ack)
}
