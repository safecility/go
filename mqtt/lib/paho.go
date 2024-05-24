package lib

import (
	"crypto/tls"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"time"
)

// Token is returned on asynchronous functions
type Token interface {
	// Wait for the function to finish
	Wait() bool
	// WaitTimeout wait for the function to finish or return false after a certain time
	WaitTimeout(time.Duration) bool
	// Error the error associated with the result of the function (nil if everything okay)
	Error() error
}

var (
	PublishQoS   byte = 0x00
	SubscribeQoS byte = 0x00
)

// PahoClient use a PahoClient to connect to MQTT
type PahoClient struct {
	opts          *paho.ClientOptions
	mqtt          paho.Client
	subscriptions map[string]paho.MessageHandler
}

// NewPahoClient creates a PahoClient
func NewPahoClient(id, username, password string, tlsConfig *tls.Config, brokers ...string) *PahoClient {

	opts := paho.NewClientOptions()

	pahoClient := &PahoClient{
		opts:          opts,
		subscriptions: make(map[string]paho.MessageHandler),
	}

	for _, broker := range brokers {
		pahoClient.opts.AddBroker(broker)
	}
	//give ourselves a unique client id
	pahoClient.opts.SetClientID(fmt.Sprintf("%s-%s", id, "-ac31-a014-ab21-1442"))
	pahoClient.opts.SetUsername(username)
	pahoClient.opts.SetPassword(password)
	pahoClient.opts.SetKeepAlive(30 * time.Second)
	pahoClient.opts.SetPingTimeout(10 * time.Second)
	pahoClient.opts.SetCleanSession(true)

	pahoClient.opts.TLSConfig = tlsConfig

	pahoClient.opts.SetDefaultPublishHandler(func(client paho.Client, msg paho.Message) {
		log.Warn().Str("message", string(msg.Payload())).Msg("unhandled message")
	})

	var reconnecting bool

	pahoClient.opts.SetConnectionLostHandler(func(client paho.Client, err error) {
		log.Err(err).Msg("mqtt: disconnected, reconnecting...")
		reconnecting = true
	})

	pahoClient.opts.SetOnConnectHandler(func(client paho.Client) {
		log.Info().Msg("mqtt: connected")
		if reconnecting {
			for topic, handler := range pahoClient.subscriptions {
				log.Debug().Str("topic", topic).Msg("mqtt: re-subscribing to topic")
				pahoClient.Subscribe(topic, handler)
			}
			reconnecting = false
		}
	})

	pahoClient.mqtt = paho.NewClient(pahoClient.opts)

	return pahoClient
}

var (
	// ConnectRetries says how many times the client should retry a failed connection
	ConnectRetries = 5
	// ConnectRetryDelay says how long the client should wait between retries
	ConnectRetryDelay = 2 * time.Second
)

// Connect to the MQTT broker. It will retry for ConnectRetries times with a delay of ConnectRetryDelay between retries
func (c *PahoClient) Connect() error {
	if c.mqtt == nil {
		return fmt.Errorf("no client to connect")
	}
	if c.mqtt.IsConnected() {
		return nil
	}
	var err error
	for retries := 0; retries < ConnectRetries; retries++ {
		token := c.mqtt.Connect()
		finished := token.WaitTimeout(1 * time.Second)
		if !finished {
			log.Warn().Msg("mqtt: connection took longer than expected...")
			token.Wait()
		}
		err = token.Error()
		if err == nil {
			break
		}
		log.Err(err).Msg("mqtt: could not connect, retrying...")
		<-time.After(ConnectRetryDelay)
	}
	if err != nil {
		return fmt.Errorf("could not connect to MQTT Broker (%s)", err)
	}
	return nil
}

func (c *PahoClient) Publish(topic string, msg []byte) Token {
	if c.mqtt == nil {
		return nil
	}
	return c.mqtt.Publish(topic, PublishQoS, false, msg)
}

func (c *PahoClient) Subscribe(topic string, handler paho.MessageHandler) Token {
	if c.mqtt == nil {
		return nil
	}
	c.subscriptions[topic] = handler
	return c.mqtt.Subscribe(topic, SubscribeQoS, handler)
}

func (c *PahoClient) unsubscribe(topic string) Token {
	if c.mqtt == nil {
		return nil
	}
	delete(c.subscriptions, topic)
	return c.mqtt.Unsubscribe(topic)
}

// Disconnect from the MQTT broker
func (c *PahoClient) Disconnect() {
	if c.mqtt == nil {
		log.Warn().Msg("no client to disconnect")
		return
	}
	if !c.mqtt.IsConnected() {
		return
	}
	log.Debug().Msg("mqtt: disconnecting")
	c.mqtt.Disconnect(25)
}

// IsConnected returns true if there is a connection to the MQTT broker
func (c *PahoClient) IsConnected() bool {
	if c.mqtt == nil {
		return false
	}
	return c.mqtt.IsConnected()
}
