package mqtt

import (
	"fmt"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	//util "gisplatform-basemap-api/util/array"
	//"time"
)

func TestStart(t *testing.T) {
	//MQTTService1.Option().SetPassword("asdsadsad")
	//_, err := MQTTService1.Start()
	//MQTTService.Option().AddBroker("tcp://127.0.0.1:1883")
	_, err := MQTTService.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo connection mqtt: %v", err)
	}
	if !MQTTService.Ping() {
		t.Fail()
	}
	MQTTService.Stop()
}

func TestStop(t *testing.T) {
	_, err := MQTTService.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo connection mqtt: %v", err)
	}
	err1 := MQTTService.Stop()
	if err1 != nil {
		t.Errorf("Lỗi khi đóng connect: %v", err1)
	}
}

func Test_Publish_1(t *testing.T) {
	_, err := MQTTService.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo connection mqtt: %v", err)
	}
	err1 := MQTTService.Publish(Value1.Topic, Value1.Message, false)
	if err1 != nil {
		t.Errorf("Lỗi publish: %v", err1)
	}

	MQTTService.Stop()
}

func Test_Publish_2(t *testing.T) {
	_, err := MQTTService.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo connection mqtt: %v", err)
	}
	err1 := MQTTService.PublishQOS(Value1.Topic, Value1.Message, false, 0)
	if err1 != nil {
		t.Errorf("Lỗi publish: %v", err1)
	}
	err2 := MQTTService.PublishQOS(Value2.Topic, Value2.Message, false, 1)
	if err2 != nil {
		t.Errorf("Lỗi publish: %v", err2)
	}

	MQTTService.Stop()
}

func Test_Publish_3(t *testing.T) {
	_, err := MQTTService.Start()

	if err != nil {
		t.Errorf("Lỗi khởi tạo connection mqtt: %v", err)
	}

	err1 := MQTTService.PublishQOS(Value1.Topic, Value1.Message, false, 0)
	if err1 != nil {
		t.Errorf("Lỗi publish: %v", err1)
	}

	err2 := MQTTService.PublishQOS(Value2.Topic, Value2.Message, false, 1)
	if err2 != nil {
		t.Errorf("Lỗi publish: %v", err2)
	}

	err3 := MQTTService.PublishQOS(Value3.Topic, Value3.Message, false, 2)
	if err3 != nil {
		t.Errorf("Lỗi publish: %v", err3)
	}

	MQTTService.Stop()
}
func Test_PublishQOS(t *testing.T) {
	_, err := MQTTService.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
	}

	err1 := MQTTService.PublishQOS(Value1.Topic, Value1.Message, false, Value1.Qos)

	if err1 != nil {
		t.Errorf("Lỗi publish: %v", err1)
	}

	MQTTService.Stop()
}

func Test_Publish_Auth(t *testing.T) {
	ops := paho.NewClientOptions().AddBroker("tcp://" + "mqtt.dioty.co" + ":1883").SetUsername("salt.coffee189@gmail.com").SetPassword("bf64fe27")
	ops.SetClientID("Publish_3")
	c := paho.NewClient(ops)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	c.Publish("/salt.coffee189@gmail.com/test", 0, false, "Publish2 qos0")
	c.Publish("/salt.coffee189@gmail.com/test", 1, false, "Publish2 qos1")
	c.Publish("/salt.coffee189@gmail.com/test", 2, false, "Publish2 qos2")

	c.Disconnect(500)
}

func Test_Subscribe(t *testing.T) {
	var f paho.MessageHandler = func(client paho.Client, msg paho.Message) {
		fmt.Printf("[%s] %s\n", msg.Topic(), string(msg.Payload()))
	}
	MQTTService1.Option().SetClientID("Subscribe_rx").SetDefaultPublishHandler(f)
	_, err := MQTTService1.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
	}

	errsub := MQTTService1.Subscribe("/salt.coffee189@gmail.com/test", nil, 0)
	if errsub != nil {
		t.Errorf("Lỗi subcribe: %v", errsub)
	}
	_, err1 := MQTTService.Start()
	if err1 != nil {
		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
	}
	errpub := MQTTService.Publish("/salt.coffee189@gmail.com/test", "Publish qos0", true)
	if errpub != nil {
		t.Errorf("Lỗi publish: %v", errpub)
	}

	MQTTService.Stop()
	MQTTService1.Stop()
}

//Terminal this function when running to see will message
// func Test_Will(t *testing.T) {
// 	willmsgc := make(chan string, 1)
// 	MQTTService1.Option().SetClientID("will-giver")
// 	MQTTService1.Option().SetWill("/salt.coffee189@gmail.com/wills", "good-byte!", 0, false)
// 	MQTTService1.Option().SetConnectionLostHandler(func(client paho.Client, err error) {
// 		fmt.Println("OnConnectionLost!")
// 	})

// 	MQTTService.Option().SetClientID("will-subscriber")
// 	MQTTService.Option().SetDefaultPublishHandler(func(client paho.Client, msg paho.Message) {
// 		fmt.Printf("TOPIC: %s\n", msg.Topic())
// 		fmt.Printf("MSG: %s\n", msg.Payload())
// 		willmsgc <- string(msg.Payload())
// 	})

// 	_, err := MQTTService.Start()
// 	if err != nil {
// 		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
// 	}

// 	errsub := MQTTService.Subscribe("/salt.coffee189@gmail.com/wills", nil, 0)
// 	if errsub != nil {
// 		t.Errorf("Lỗi subcribe: %v", errsub)
// 	}

// 	_, err1 := MQTTService1.Start()
// 	if err1 != nil {
// 		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
// 	}

// 	if <-willmsgc != "good-byte!" {
// 		t.Fatalf("will message did not have correct payload")
// 	}

// 	time.Sleep(100 * time.Second)
// }

// func wait(c chan bool) {
// 	fmt.Println("choke is waiting")
// 	<-c
// }

func TestUnSubscribe(t *testing.T) {
	_, err := MQTTService.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
	}
	err1 := MQTTService.UnSubscribe(Value1.Topic)
	if err1 != nil {
		t.Errorf("Lỗi unsubcribe")
	}
}

//Turn off your network in 5 second to let the mqtt fall to reconnect state then turn on back your network
func Test_Reconnect(t *testing.T) {
	MQTTService.Option().SetClientID("auto_reconnect")
	MQTTService.Option().SetKeepAlive(time.Duration(2) * time.Second)
	MQTTService.Option().SetPingTimeout(time.Duration(1) * time.Second)
	_, err := MQTTService.Start()
	if err != nil {
		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
	}

	//fmt.Println("Breaking connection in 5 second then reconnect")
	time.Sleep(15 * time.Second)

	if !MQTTService.Ping() {
		t.Fail()
	}

	MQTTService.Stop()
}

//Test with local broker
// func Test_PublishMessage(t *testing.T) {
// 	topic := "pubmsg"
// 	choke := make(chan bool)
// 	MQTTService.Option().SetClientID("pubmsg-pub")

// 	MQTTService1.Option().SetClientID("pubmsg-sub")
// 	var f paho.MessageHandler = func(client paho.Client, msg paho.Message) {
// 		fmt.Printf("TOPIC: %s\n", msg.Topic())
// 		fmt.Printf("MSG: %s\n", msg.Payload())
// 		if string(msg.Payload()) != "pubmsg payload" {
// 			fmt.Println("Message payload incorrect", msg.Payload(), len("pubmsg payload"))
// 			t.Fatalf("Message payload incorrect")
// 		}
// 		choke <- true
// 	}
// 	MQTTService1.Option().SetDefaultPublishHandler(f)

// 	_, err := MQTTService1.Start()
// 	if err != nil {
// 		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
// 	}
// 	MQTTService1.Subscribe(topic, nil, 2)
// 	_, err2 := MQTTService.Start()
// 	if err2 != nil {
// 		t.Errorf("Lỗi khởi tạo mqtt: %v", err)
// 	}
// 	mqtt.
// 	text := "pubmsg payload"
// 	MQTTService.Publish(topic, text, false)
// 	MQTTService.Publish(topic, text, false)

// 	wait(choke)
// 	wait(choke)

// 	MQTTService.Publish(topic, text, false)
// 	wait(choke)

// 	MQTTService.Stop()
// 	MQTTService1.Stop()
// }
