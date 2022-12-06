package mqttTooling

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"reflect"
	"testing"
)

func TestTopicFormatter_Apply(t *testing.T) {
	type args struct {
		room  string
		topic string
	}
	tests := []struct {
		name string
		t    TopicFormatter
		args args
		want Topic
	}{
		{name: "default", t: TopicFormatter("room/%s/topic/%s"), args: args{"bedroom", "state"}, want: "room/bedroom/topic/state"},
		{name: "setter", t: TopicFormatter("room/%s/topic/%s/set"), args: args{"bedroom", "state"}, want: "room/bedroom/topic/state/set"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Apply(tt.args.room, tt.args.topic); got != tt.want {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMqttPublisher_Publish(t *testing.T) {
	client := ClientMock{
		PublishChan: make(chan MqttMsg, 1),
	}
	defer close(client.PublishChan)

	type fields struct {
		client mqtt.Client
	}
	type args struct {
		topic   Topic
		payload []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "simple msg",
			fields: fields{client: &client},
			args: args{
				topic:   "simple/msg",
				payload: []byte("content msg"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MqttPublisher{
				client: tt.fields.client,
			}
			err := m.Publish(tt.args.topic, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				expectedMsg := MqttMsg{Topic: tt.args.topic, Qos: byte(0), Payload: tt.args.payload}
				msg := <-client.PublishChan
				if !reflect.DeepEqual(msg, expectedMsg) {
					t.Errorf("bad message published: '%v', want '%v'", msg, expectedMsg)
				}
			}
		})
	}
}

type ClientMock struct {
	PublishChan chan MqttMsg
}

func (c ClientMock) IsConnected() bool {
	panic("implement me")
}

func (c ClientMock) IsConnectionOpen() bool {
	panic("implement me")
}

func (c ClientMock) Connect() mqtt.Token {
	panic("implement me")
}

func (c ClientMock) Disconnect(quiesce uint) {
	panic("implement me")
}

func (c ClientMock) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	c.PublishChan <- MqttMsg{Topic(topic), qos, retained, payload}
	return &mqtt.DummyToken{}
}

func (c ClientMock) PublishAsString(topic Topic, qos byte, retained bool, payload fmt.Stringer) mqtt.Token {
	c.PublishChan <- MqttMsg{topic, qos, retained, []byte(payload.String())}
	return &mqtt.DummyToken{}
}

func (c ClientMock) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	panic("implement me")
}

func (c ClientMock) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	panic("implement me")
}

func (c ClientMock) Unsubscribe(topics ...string) mqtt.Token {
	panic("implement me")
}

func (c ClientMock) AddRoute(topic string, callback mqtt.MessageHandler) {
	panic("implement me")
}

func (c ClientMock) OptionsReader() mqtt.ClientOptionsReader {
	panic("implement me")
}

type MqttMsg struct {
	Topic    Topic
	Qos      byte
	Retained bool
	Payload  interface{}
}
