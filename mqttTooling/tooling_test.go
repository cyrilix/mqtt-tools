package mqttTooling

import (
	"github.com/cyrilix/mqtt-tools/mqttTooling/mqtttest"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"reflect"
	"testing"
)

func TestMqttPublisher_Publish(t *testing.T) {
	client := mqtttest.ClientMock{
		PublishChan: make(chan mqtttest.MqttMsg, 1),
	}
	defer close(client.PublishChan)

	type fields struct {
		client mqtt.Client
	}
	type args struct {
		topic   string
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
				expectedMsg := mqtttest.MqttMsg{Topic: tt.args.topic, Qos: byte(0), Payload: tt.args.payload}
				msg := <-client.PublishChan
				if !reflect.DeepEqual(msg, expectedMsg) {
					t.Errorf("bad message published: '%v', want '%v'", msg, expectedMsg)
				}
			}
		})
	}
}
