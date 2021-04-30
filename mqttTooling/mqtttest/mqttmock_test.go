package mqtttest

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"reflect"
	"testing"
)

func TestClientMock_Publish(t *testing.T) {
	type fields struct {
		PublishChan chan MqttMsg
	}
	type args struct {
		topic    string
		qos      byte
		retained bool
		payload  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   mqtt.Token
	}{
		{
			"simple test",
			fields{
				PublishChan: make(chan MqttMsg, 1),
			},
			args{
				topic:    "topicTest",
				qos:      0,
				retained: false,
				payload:  "msg content",
			},
			&mqtt.DummyToken{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ClientMock{
				PublishChan: tt.fields.PublishChan,
			}
			if got := c.Publish(tt.args.topic, tt.args.qos, tt.args.retained, tt.args.payload); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Publish() = %v, want %v", got, tt.want)
			}
			msgPublished := <-tt.fields.PublishChan
			if msgPublished.Topic != tt.args.topic {
				t.Errorf("msg published to bad topic '%v', want '%v'", msgPublished.Topic, tt.args.topic)
			}
			if msgPublished.Qos != tt.args.qos {
				t.Errorf("msg published with bad qos '%v', want '%v'", msgPublished.Qos, tt.args.qos)
			}
			if msgPublished.Retained != tt.args.retained {
				t.Errorf("msg published with bad retained flag '%v', want '%v'", msgPublished.Retained, tt.args.retained)
			}
			if msgPublished.Payload != tt.args.payload {
				t.Errorf("msg published with bad payload '%v', want '%v'", msgPublished.Payload, tt.args.payload)
			}
			close(tt.fields.PublishChan)
		})
	}
}
