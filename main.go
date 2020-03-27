package main

import (
	"encoding/json"
	"time"

	MQTT "bitbucket.org/mqttgis/transport"
)

var knt int

// LatLng - struct
type LatLng struct {
	LAT float64 `json:"lat"`
	LON float64 `json:"lon"`
}

func main() {
	// knt = 0
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	client := MQTT.CreateMQTTClient()
	client.Start()
	// opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt.dioty.co:1883")
	// opts.SetClientID("phuc123")
	// opts.SetUsername("salt.coffee189@gmail.com")
	// opts.SetPassword("bf64fe27")
	// opts.SetDefaultPublishHandler(f)
	// topic := "phuc-topic"

	// opts.OnConnect = func(c mqtt.Client) {
	// 	if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
	// 		panic(token.Error())
	// 	}
	// }

	var GPS [10]LatLng
	// GPS[0] = LatLng{LAT: 10.502746124283656, LON: 107.38368597037162}
	// GPS[1] = LatLng{LAT: 10.553810837476034, LON: 107.20132184125767}
	// GPS[2] = LatLng{LAT: 10.548570829582767, LON: 107.17570386943207}
	// GPS[3] = LatLng{LAT: 10.458918206895515, LON: 107.07248052315317}
	// GPS[4] = LatLng{LAT: 10.544808701831574, LON: 107.01121969501403}
	// GPS[5] = LatLng{LAT: 10.4491585074751001, LON: 107.23079953359036}
	// GPS[6] = LatLng{LAT: 10.468135281158133, LON: 107.00127921464853}
	// GPS[7] = LatLng{LAT: 10.450887946734518, LON: 107.35589569018214}
	// GPS[8] = LatLng{LAT: 10.490830999209736, LON: 107.0104126888608}
	// GPS[9] = LatLng{LAT: 10.472714857173292, LON: 107.24762576346343}
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

	// var f = func(c MQTT.Client, msg MQTT.Message) {
	// 	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	// }
	// client := mqtt.NewClient(opts)
	for i := 0; i < len(GPS); i++ {
		data, _ := json.Marshal(&GPS[i])
		client.Publish("/salt.coffee189@gmail.com/test", string(data))
		// client.SubscribeDefault("/salt.coffee189@gmail.com/test")
		// client.Subscribe("/salt.coffee189@gmail.com/test", f, 0)
		time.Sleep(3 * time.Second)
	}
}
