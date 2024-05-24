package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/mqtt/lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Client interface for a paho server
type Client interface {
	// Close the client and clean up all connections
	Close() error
	// GetPahoPubSub PubSub Subscribe to uplink and events, publish downlink
	GetPahoPubSub() (*PahoPubSub, error)
}

type ClientConfig struct {
	ClientName     string
	ClientVersion  string
	TLSConfig      *tls.Config
	RequestTimeout time.Duration

	appID        string
	userName     string
	appAccessKey string
	mqttAddress  string
}

// NewClient creates a new API client from the configuration, using the given Application ID and Application access key.
func NewClient(appID, userName, appAccessKey, mqttAddress string) Client {

	config := ClientConfig{
		ClientName:     "safecility",
		ClientVersion:  "1.0.0",
		RequestTimeout: 10 * time.Second,
		appID:          appID,
		userName:       userName,
		appAccessKey:   appAccessKey,
		mqttAddress:    mqttAddress,
	}

	client := &client{
		ClientConfig:         config,
		transportCredentials: credentials.NewTLS(config.TLSConfig),
	}

	return client
}

type client struct {
	ClientConfig ClientConfig

	transportCredentials credentials.TransportCredentials

	handler struct {
		sync.RWMutex
		conn *grpc.ClientConn
	}
	mqtt struct {
		sync.RWMutex
		client *lib.PahoClient
		ctx    context.Context
		cancel context.CancelFunc
	}
}

func (c *client) GetPahoPubSub() (*PahoPubSub, error) {
	if err := c.connectMQTT(); err != nil {
		return nil, err
	}

	errors := make(chan error)

	pubSub := &PahoPubSub{
		client:  c.mqtt.client,
		SubPubs: make(map[string]string, 10),
		Errors:  errors,
	}
	return pubSub, nil
}

func (c *client) Close() (closeErr error) {
	if err := c.closeHandler(); err != nil {
		closeErr = err
	}

	if err := c.closeMQTT(); err != nil && closeErr == nil {
		closeErr = err
	}
	return
}

func (c *client) closeHandler() error {
	c.handler.Lock()
	defer c.handler.Unlock()
	if c.handler.conn != nil {
		go func() {
			err := c.handler.conn.Close()
			if err != nil {
				log.Err(err).Msg("could not close handler")
			}
		}()
	}
	c.handler.conn = nil
	return nil
}

var mqttBufferSize = 10

func (c *client) connectMQTT() (err error) {
	c.mqtt.Lock()
	defer c.mqtt.Unlock()
	if c.mqtt.client != nil {
		return nil
	}
	c.handler.RLock()
	defer c.handler.RUnlock()
	if c.ClientConfig.mqttAddress == "" {
		//err = c.discoverAddress()
		return fmt.Errorf("no mqttAddress")
	}
	mqttAddress, err := cleanMQTTAddress(c.ClientConfig.mqttAddress)
	if err != nil {
		return err
	}

	var tlsConfig *tls.Config
	if strings.HasPrefix(mqttAddress, "ssl://") {
		tlsConfig = lib.NewTlsConfig()
	}
	c.mqtt.client = lib.NewPahoClient(
		c.ClientConfig.ClientName, c.ClientConfig.userName, c.ClientConfig.appAccessKey, tlsConfig, mqttAddress)

	c.mqtt.ctx, c.mqtt.cancel = context.WithCancel(context.Background())
	log.Debug().Str("Address", mqttAddress).Msg("Connecting to MQTT...")

	if err := c.mqtt.client.Connect(); err != nil {
		return fmt.Errorf("could not connect to MQTT %v", err)
	}
	log.Debug().Msg("Connected to MQTT")
	return nil
}

func (c *client) closeMQTT() error {
	c.mqtt.Lock()
	defer c.mqtt.Unlock()
	if c.mqtt.client == nil {
		return nil
	}
	log.Debug().Msg("disconnecting from MQTT...")
	c.mqtt.cancel()
	c.mqtt.client.Disconnect()
	c.mqtt.client = nil
	return nil
}

func cleanMQTTAddress(in string) (address string, err error) {
	if !strings.Contains(in, "://") {
		in = "detect://" + in
	}
	uri, err := url.Parse(in)
	if err != nil {
		return address, err
	}
	switch uri.Scheme {
	case "detect", "mqtt", "mqtts", "ssl":
	default:
		return address, fmt.Errorf("ttn-sdk: unknown mqtt scheme: %s", uri.Scheme)
	}
	scheme, host, port := uri.Scheme, uri.Hostname(), uri.Port()
	if scheme == "detect" {
		switch port {
		case "8883", "":
			scheme = "ssl"
		default:
			scheme = "tcp"
		}
	}
	if port == "" {
		switch scheme {
		case "ssl", "mqtts":
			scheme = "ssl"
			port = "8883"
		case "tcp", "mqtt":
			scheme = "tcp"
			port = "1883"
		}
	}
	return fmt.Sprintf("%s://%s:%s", scheme, host, port), nil
}
