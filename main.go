package main

import (
	"fmt"
	"time"

	IM "bitbucket.org/mqttgis/mqtt"
	mqttcon "bitbucket.org/mqttgis/mqtt/driver"
	paho "github.com/eclipse/paho.mqtt.golang"
)

var knt int

//MQTTService - khai bao MQTT lib
var MQTTService IM.IMQTTClient

//MQ - khai bao MQTT lib
var MQ IM.IMQTTClient

// LatLng - struct
type LatLng struct {
	LAT float64 `json:"lat"`
	LON float64 `json:"lon"`
}

//
func onMessageReceived(client paho.Client, message paho.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
}

//Khởi tạo
func init() {
	MQ = mqttcon.Init("./config/config.json")
	// MQTTService = mqttcon.InitMQTTClientOptions("./config/config.json")
}

func main() {
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//MQTTService.SubscribeOC("topic", onMessageReceived, 0)
	// MQTTService.CreateMQTTClient()
	// MQTTService.Start()
	MQ.SubscribeD("/salt.coffee189@gmail.com/huydo189")
	//MQTTService.SubscribeD("/salt.coffee189@gmail.com/huydo189")
	var GPS [10]LatLng
	GPS[0] = LatLng{LAT: 10.445595, LON: 107.186590}
	GPS[1] = LatLng{LAT: 10.510079, LON: 107.246316}
	GPS[2] = LatLng{LAT: 10.512941, LON: 107.276870}
	GPS[3] = LatLng{LAT: 10.556337, LON: 107.321678}
	GPS[4] = LatLng{LAT: 10.572473, LON: 107.274554}
	GPS[5] = LatLng{LAT: 10.562583, LON: 107.231137}
	GPS[6] = LatLng{LAT: 10.521461, LON: 107.244374}
	GPS[7] = LatLng{LAT: 10.496993, LON: 107.277607}
	GPS[8] = LatLng{LAT: 10.474085, LON: 107.248610}
	GPS[9] = LatLng{LAT: 10.477209, LON: 107.297322}
	for i := 0; i < 100; i++ {
		// data, _ := json.Marshal(&GPS[i])
		// MQTTService.Publish("/salt.coffee189@gmail.com/huydo189", string(data))
		time.Sleep(3 * time.Second)
	}
}
