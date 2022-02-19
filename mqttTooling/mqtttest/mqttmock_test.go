package mqtttest

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
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

func TestPublisherMock_PublishSubscribe(t *testing.T) {
	subscribeResponse := make(chan []byte)
	type fields struct {
	}
	type args struct {
		topic string
		mh    mqtt.MessageHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "default",
			fields: fields{},
			args: args{
				topic: "topic1",
				mh: func(client mqtt.Client, message mqtt.Message) {
					log.Info("OK")
					subscribeResponse <- message.Payload()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPubSub()
			p.Subscribe(tt.args.topic, tt.args.mh)
			payloadPublished := []byte(fmt.Sprintf("message test for topic %s", tt.args.topic))
			err := p.Publish(tt.args.topic, payloadPublished)
			if err != nil {
				t.Errorf("unexepected error: %v", err)
			}
			select {
			case payload := <-subscribeResponse:
				if string(payload) != string(payloadPublished) {
					t.Errorf("invalid message received: %s, want %s", string(payload), string(payloadPublished))
				}
			case <-time.NewTimer(100 * time.Millisecond).C:
				t.Errorf("no msg received")
			}
		})
	}
}
