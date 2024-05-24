package stream

import (
	"cloud.google.com/go/pubsub"
	"reflect"
	"testing"
)

func TestPublishToTopic(t *testing.T) {
	type args struct {
		m     interface{}
		topic *pubsub.Topic
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestPublishToTopic nil topic",
			args: args{
				m:     nil,
				topic: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublishToTopic(tt.args.m, tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublishToTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublishToTopic() got = %v, want %v", got, tt.want)
			}
		})
	}
}
