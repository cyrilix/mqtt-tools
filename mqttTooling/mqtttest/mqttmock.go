package mqtttest

import (
	"github.com/cyrilix/mqtt-tools/mqttTooling"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewPubSub() mqttTooling.MQTTPubSub {
	return &PubSubMock{make(map[string]chan MqttMsg, 1), make(chan interface{})}
}

type PubSubMock struct {
	topicChan map[string]chan MqttMsg
	cancel    chan interface{}
}

func (p *PubSubMock) Subscribe(topic string, mh mqtt.MessageHandler) {
	if _, ok := p.topicChan[topic]; !ok {
		p.topicChan[topic] = make(chan MqttMsg)
	}
	go func() {
		var msg MqttMsg
		for {
			select {
			case msg = <-p.topicChan[topic]:
				mh(nil, &MessageMock{payload: msg.Payload.([]byte), topic: msg.Topic, qos: msg.Qos, retained: msg.Retained})
			case <-p.cancel:
				return
			}
		}
	}()
}

func (p *PubSubMock) Close() error {
	p.cancel <- struct{}{}
	return nil
}

func (p *PubSubMock) Publish(topic string, payload []byte) error {
	if _, ok := p.topicChan[topic]; !ok {
		p.topicChan[topic] = make(chan MqttMsg)
	}
	p.topicChan[topic] <- MqttMsg{
		Topic:    topic,
		Qos:      0,
		Retained: false,
		Payload:  payload,
	}
	return nil
}

type MqttMsg struct {
	Topic    string
	Qos      byte
	Retained bool
	Payload  interface{}
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
	c.PublishChan <- MqttMsg{topic, qos, retained, payload}
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

type MessageMock struct {
	duplicate bool
	qos       byte
	retained  bool
	topic     string
	messageId uint16
	payload   []byte
}

func (m *MessageMock) Duplicate() bool {
	return m.duplicate
}

func (m *MessageMock) Qos() byte {
	return m.qos
}

func (m *MessageMock) Retained() bool {
	return m.retained
}

func (m *MessageMock) Topic() string {
	return m.topic
}

func (m *MessageMock) MessageID() uint16 {
	return m.messageId
}

func (m *MessageMock) Payload() []byte {
	return m.payload
}

func (m *MessageMock) Ack() {
}
