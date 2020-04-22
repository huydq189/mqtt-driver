package mqtt

import (
	mqttconn "bitbucket.org/mqttgis/mqtt/driver"
)

//Test - struct
type Test struct {
	Topic   string
	Qos     byte
	Message string
}

var pathconfig = `../../../config/config.json`

//MQTTService - test MQService
var MQTTService = mqttconn.InitMQTTClientOptions(pathconfig)

//MQ - test MQService
var MQ = mqttconn.Init(pathconfig)

//Value - test
var Value = Test{
	Topic:   "/salt.coffee189@gmail.com/test",
	Qos:     0,
	Message: "Hello from the other side",
}
