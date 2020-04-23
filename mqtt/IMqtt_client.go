package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// IMQTTClient - IModelRepository
type IMQTTClient interface {
	Start() (bool, error)
	Stop() error
	SetOnConnectHandler(handler mqtt.OnConnectHandler) *mqtt.ClientOptions
	Subscribe(topic string, handler mqtt.MessageHandler, qos byte) error
	//SubscribeOC(topic string, handler mqtt.MessageHandler, qos byte) *mqtt.ClientOptions
	SubscribeD(topic string) error
	UnSubscribe(topic string) error
	Publish(topic, message string, retained bool) error
	PublishQOS(topic, message string, retained bool, qos byte) error
}
