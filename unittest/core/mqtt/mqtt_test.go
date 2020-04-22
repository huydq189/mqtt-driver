package mqtt

import (
	"testing"

	paho "github.com/eclipse/paho.mqtt.golang"
	//util "gisplatform-basemap-api/util/array"
	//"time"
)

func TestCreateMQTTClient(t *testing.T) {
	err := MQTTService.CreateMQTTClient()
	if err != nil {
		t.Errorf("Lỗi tạo client và clientoptions mqtt: %v", err)
	}
}

func TestStart(t *testing.T) {
	err := MQTTService.CreateMQTTClient()
	if err != nil {
		t.Errorf("Lỗi tạo client và clientoptions mqtt: %v", err)
	}
	_, err1 := MQTTService.Start()
	if err1 != nil {
		t.Errorf("Lỗi connect broker: %v", err)
	}
}

func TestStop(t *testing.T) {
	err := MQ.Stop()
	if err != nil {
		t.Errorf("Lỗi khi đóng connect: %v", err)
	}
}

func TestSubscribe(t *testing.T) {
	f := func(client paho.Client, message paho.Message) {
		if string(message.Payload()) != Value.Message {
			t.Errorf("Message nhận được không đúng")
		}
	}
	err := MQ.Subscribe(Value.Topic, f, Value.Qos)
	if err != nil {
		t.Errorf("Lỗi subcribe: %v", err)
	}
	MQ.Publish(Value.Topic, Value.Message)
}

func TestUnSubscribe(t *testing.T) {
	err := MQ.UnSubscribe(Value.Topic)
	if err != nil {
		t.Errorf("Lỗi unsubcribe")
	}
}

func TestPublish(t *testing.T) {
	f := func(client paho.Client, message paho.Message) {
		if string(message.Payload()) != Value.Message {
			t.Errorf("Message nhận được không đúng")
		}
	}
	err := MQ.Subscribe(Value.Topic, f, Value.Qos)
	if err != nil {
		t.Errorf("Lỗi subcribe: %v", err)
	}
	err1 := MQ.Publish(Value.Topic, Value.Message)
	if err1 != nil {
		t.Errorf("Lỗi publish: %v", err)
	}
}

func TestPublishQOS(t *testing.T) {
	f := func(client paho.Client, message paho.Message) {
		if string(message.Payload()) != Value.Message {
			t.Errorf("Message nhận được không đúng")
		}
	}
	err := MQ.Subscribe(Value.Topic, f, Value.Qos)
	if err != nil {
		t.Errorf("Lỗi subcribe: %v", err)
	}
	err1 := MQ.PublishQOS(Value.Topic, Value.Message, Value.Qos)
	if err1 != nil {
		t.Errorf("Lỗi publish: %v", err)
	}
}
