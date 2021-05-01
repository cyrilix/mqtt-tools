package mqtttest

import mqtt "github.com/eclipse/paho.mqtt.golang"

func NewPublisherMock() *PublisherMock {
	return &PublisherMock{make(chan MqttMsg, 1)}
}

type PublisherMock struct {
	PublishChan chan MqttMsg
}

func (p PublisherMock) Close() error {
	close(p.PublishChan)
	return nil
}

func (p PublisherMock) Publish(topic string, payload []byte) error {
	p.PublishChan <- MqttMsg{
		Topic:   topic,
		Payload: payload,
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
