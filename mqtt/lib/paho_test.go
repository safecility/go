package lib

import (
	"crypto/tls"
	paho "github.com/eclipse/paho.mqtt.golang"
	"reflect"
	"testing"
	"time"
)

func TestNewPahoClient(t *testing.T) {
	type args struct {
		id        string
		username  string
		password  string
		tlsConfig *tls.Config
		brokers   []string
	}
	tests := []struct {
		name string
		args args
		want *PahoClient
	}{
		// TODO: Add test cases.
		{
			name: "setup",
			args: args{
				id:        "",
				username:  "example",
				password:  "test",
				tlsConfig: nil,
				brokers:   nil,
			},
			want: &PahoClient{
				opts: &paho.ClientOptions{
					Servers:                nil,
					ClientID:               "",
					Username:               "example",
					Password:               "test",
					CredentialsProvider:    nil,
					CleanSession:           true,
					Order:                  false,
					WillEnabled:            false,
					WillTopic:              "",
					WillPayload:            nil,
					WillQos:                0,
					WillRetained:           false,
					ProtocolVersion:        0,
					TLSConfig:              nil,
					KeepAlive:              0,
					PingTimeout:            10 * time.Second,
					ConnectTimeout:         0,
					MaxReconnectInterval:   0,
					AutoReconnect:          false,
					ConnectRetryInterval:   0,
					ConnectRetry:           false,
					Store:                  nil,
					DefaultPublishHandler:  nil,
					OnConnect:              nil,
					OnConnectionLost:       nil,
					OnReconnecting:         nil,
					OnConnectAttempt:       nil,
					WriteTimeout:           0,
					MessageChannelDepth:    0,
					ResumeSubs:             false,
					HTTPHeaders:            nil,
					WebsocketOptions:       nil,
					MaxResumePubInFlight:   0,
					Dialer:                 nil,
					CustomOpenConnectionFn: nil,
					AutoAckDisabled:        false,
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := NewPahoClient(tt.args.id, tt.args.username, tt.args.password, tt.args.tlsConfig, tt.args.brokers...); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewPahoClient() = %v, want %v", got, tt.want)
			//}
			got := NewPahoClient(tt.args.id, tt.args.username, tt.args.password, tt.args.tlsConfig)
			if got.opts.Username != tt.want.opts.Username {
				t.Errorf("NewPahoClient() username = %v, want %v", got.opts.Username, tt.want.opts.Username)
			}
			if got.opts.Password != tt.want.opts.Password {
				t.Errorf("NewPahoClient() password = %v, want %v", got.opts.Password, tt.want.opts.Password)
			}
			if got.opts.CleanSession != tt.want.opts.CleanSession {
				t.Errorf("NewPahoClient() CleanSession = %v, want %v", got.opts.CleanSession, tt.want.opts.CleanSession)
			}
		})
	}
}

func TestPahoClient_Connect(t *testing.T) {
	type fields struct {
		opts          *paho.ClientOptions
		mqtt          paho.Client
		subscriptions map[string]paho.MessageHandler
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add more/better test cases.
		{
			name: "no mqtt client",
			fields: fields{
				opts:          &paho.ClientOptions{},
				mqtt:          nil,
				subscriptions: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PahoClient{
				opts:          tt.fields.opts,
				mqtt:          tt.fields.mqtt,
				subscriptions: tt.fields.subscriptions,
			}
			if err := c.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPahoClient_Disconnect(t *testing.T) {
	type fields struct {
		opts          *paho.ClientOptions
		mqtt          paho.Client
		subscriptions map[string]paho.MessageHandler
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add more/better test cases.
		{
			name: "no connection",
			fields: fields{
				opts:          &paho.ClientOptions{},
				mqtt:          nil,
				subscriptions: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PahoClient{
				opts:          tt.fields.opts,
				mqtt:          tt.fields.mqtt,
				subscriptions: tt.fields.subscriptions,
			}
			c.Disconnect()
		})
	}
}

func TestPahoClient_IsConnected(t *testing.T) {
	type fields struct {
		opts          *paho.ClientOptions
		mqtt          paho.Client
		subscriptions map[string]paho.MessageHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add more/better test cases.
		{
			name: "no connection",
			fields: fields{
				opts:          &paho.ClientOptions{},
				mqtt:          nil,
				subscriptions: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PahoClient{
				opts:          tt.fields.opts,
				mqtt:          tt.fields.mqtt,
				subscriptions: tt.fields.subscriptions,
			}
			if got := c.IsConnected(); got != tt.want {
				t.Errorf("IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPahoClient_Publish(t *testing.T) {
	type fields struct {
		opts          *paho.ClientOptions
		mqtt          paho.Client
		subscriptions map[string]paho.MessageHandler
	}
	type args struct {
		topic string
		msg   []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PahoClient{
				opts:          tt.fields.opts,
				mqtt:          tt.fields.mqtt,
				subscriptions: tt.fields.subscriptions,
			}
			if got := c.Publish(tt.args.topic, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Publish() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPahoClient_Subscribe(t *testing.T) {
	type fields struct {
		opts          *paho.ClientOptions
		mqtt          paho.Client
		subscriptions map[string]paho.MessageHandler
	}
	type args struct {
		topic   string
		handler paho.MessageHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PahoClient{
				opts:          tt.fields.opts,
				mqtt:          tt.fields.mqtt,
				subscriptions: tt.fields.subscriptions,
			}
			if got := c.Subscribe(tt.args.topic, tt.args.handler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscribe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPahoClient_unsubscribe(t *testing.T) {
	type fields struct {
		opts          *paho.ClientOptions
		mqtt          paho.Client
		subscriptions map[string]paho.MessageHandler
	}
	type args struct {
		topic string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PahoClient{
				opts:          tt.fields.opts,
				mqtt:          tt.fields.mqtt,
				subscriptions: tt.fields.subscriptions,
			}
			if got := c.unsubscribe(tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unsubscribe() = %v, want %v", got, tt.want)
			}
		})
	}
}
