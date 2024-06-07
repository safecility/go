package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/mqtt/messages"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestTtnV3_CreateDownlink(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		message        stream.SimpleMessage
		correlationIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DownlinksV3
		wantErr bool
	}{
		// TODO: Add more/better test cases.
		{
			name: "TestTtnV3_CreateDownlink",
			fields: fields{
				AppID:          "testApp",
				UidTransformer: nil,
			},
			args: args{
				message: stream.SimpleMessage{
					Payload: []byte{0, 2},
					Time:    time.Time{},
				},
				correlationIDs: nil,
			},
			want: &DownlinksV3{Downlinks: []DownlinkV3Body{
				{
					FPort:          15,
					Payload:        base64.StdEncoding.EncodeToString([]byte{0, 2}),
					Priority:       "HIGH",
					Confirmed:      true,
					CorrelationIDs: nil,
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			got := t.createV3Downlink(tt.args.message, tt.args.correlationIDs)
			log.Debug().Str("got", fmt.Sprintf("%v", got)).Msg("downlink")
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("CreateDownlink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_GetDownlinkTopicPush(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		devID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "TestTtnV3_GetDownlinkTopicPush",
			fields: fields{"testApp", nil},
			args:   args{"testDevID"},
			want:   fmt.Sprintf("v3/%s@%s/devices/%s/down/push", "testApp", "ttn", "testDevID"),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			if got := t.GetDownlinkTopicPush(tt.args.devID); got != tt.want {
				t1.Errorf("GetDownlinkTopicPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_GetDownlinkTopicReplace(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		devID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add more/better test cases.
		{
			name:   "TestTtnV3_GetDownlinkTopicReplace",
			fields: fields{"testApp", nil},
			args:   args{"testDevID"},
			want:   fmt.Sprintf("v3/%s@%s/devices/%s/down/replace", "testApp", "ttn", "testDevID"),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			if got := t.GetDownlinkTopicReplace(tt.args.devID); got != tt.want {
				t1.Errorf("GetDownlinkTopicReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_MqttTopic(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		deviceID string
		channel  messages.MqttPath
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add more/better test cases.
		{
			name:   "TestTtnV3_MqttTopic",
			fields: fields{"testApp", nil},
			args:   args{"testDevID", "mqtt"},
			want:   fmt.Sprintf("v3/%s@%s/devices/%s/%s", "testApp", "ttn", "testDevID", "mqtt"),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			if got := t.MqttTopic(tt.args.deviceID, tt.args.channel); got != tt.want {
				t1.Errorf("MqttTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_TransformPahoDownlinkMessage(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		m       mqtt.Message
		channel messages.MqttPath
	}
	pl := base64.StdEncoding.EncodeToString([]byte{222, 222})
	mID := uint16(12)
	dm := DownlinkV3{
		AppID:      "",
		MessageIDs: MessageIDs{},
		DownlinkSent: DownlinkV3Body{
			Payload: pl,
		},
	}
	js, err := json.Marshal(dm)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *stream.SimpleMessage
		wantErr bool
	}{
		// TODO: Add more/better test cases.
		{
			name:   "TestTtnV3_TransformPahoDownlinkMessage",
			fields: fields{"testApp", nil},
			args: args{&mockMessage{
				messageID: mID,
				payload:   js,
				once:      sync.Once{},
			}, messages.Sent},
			want: &stream.SimpleMessage{
				Payload: []byte{222, 222},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			got, err := t.TransformPahoDownlinkMessage(tt.args.m, tt.args.channel)
			if (err != nil) != tt.wantErr {
				t1.Errorf("TransformPahoDownlinkMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Payload, tt.want.Payload) {
				t1.Errorf("TransformPahoDownlinkMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_TransformPahoJoinMessage(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		m mqtt.Message
	}
	pl := base64.StdEncoding.EncodeToString([]byte{222, 222})
	mID := uint16(12)
	dm := JoinV3{
		AppID:      "",
		MessageIDs: MessageIDs{},
		Payload:    pl,
	}
	js, err := json.Marshal(dm)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *stream.SimpleMessage
		wantErr bool
	}{
		// TODO: Add more/better test cases.
		{
			name:   "TestTtnV3_TransformPahoJoinMessage",
			fields: fields{"testApp", nil},
			args: args{&mockMessage{
				messageID: mID,
				payload:   js,
				once:      sync.Once{},
			}},
			want: &stream.SimpleMessage{
				Payload: []byte{222, 222},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			got, err := t.TransformPahoJoinMessage(tt.args.m)
			if (err != nil) != tt.wantErr {
				t1.Errorf("TransformPahoJoinMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Payload, tt.want.Payload) {
				t1.Errorf("TransformPahoJoinMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_TransformPahoUplinkMessage(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		m mqtt.Message
	}
	pl := base64.StdEncoding.EncodeToString([]byte{222, 222})
	mID := uint16(12)
	dm := UplinkV3{
		AppID:      "",
		MessageIDs: MessageIDs{},
		Message: UplinkMessage{
			Payload:        pl,
			DecodedPayload: "aa",
			Metadata:       nil,
		},
	}
	js, err := json.Marshal(dm)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *stream.SimpleMessage
		wantErr bool
	}{
		// TODO: Add more/better test cases.
		{
			name:   "TestTtnV3_TransformPahoUplinkMessage",
			fields: fields{"testApp", nil},
			args: args{&mockMessage{
				messageID: mID,
				payload:   js,
				once:      sync.Once{},
			}},
			want: &stream.SimpleMessage{
				Payload: []byte{222, 222},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			got, err := t.TransformPahoUplinkMessage(tt.args.m)
			if (err != nil) != tt.wantErr {
				t1.Errorf("TransformPahoUplinkMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Payload, tt.want.Payload) {
				t1.Errorf("TransformPahoUplinkMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_UplinkErrorsTopic(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		deviceID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			if got := t.UplinkErrorsTopic(tt.args.deviceID); got != tt.want {
				t1.Errorf("UplinkErrorsTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTtnV3_UplinkTopic(t1 *testing.T) {
	type fields struct {
		AppID          string
		UidTransformer UidTransformer
	}
	type args struct {
		deviceID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add more/better test cases.
		{
			name:   "TestTtnV3_MqttTopic",
			fields: fields{"testApp", nil},
			args:   args{"testDevID"},
			want:   fmt.Sprintf("v3/%s@%s/devices/%s/%s", "testApp", "ttn", "testDevID", "up"),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TtnV3{
				AppID:          tt.fields.AppID,
				UidTransformer: tt.fields.UidTransformer,
			}
			if got := t.UplinkTopic(tt.args.deviceID); got != tt.want {
				t1.Errorf("UplinkTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
