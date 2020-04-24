package mqtt

import (
	mqttservice "bitbucket.org/mqttgis/core/mqtt"
)

//Test - struct
type Test struct {
	Topic   string
	Qos     byte
	Message string
}

var pathconfig = `../../../config/client.json`

//MQTTService - test MQService
var MQTTService = mqttservice.InitMQTTClientOptions(pathconfig)

//MQTTService1 - test MQService
var MQTTService1 = mqttservice.InitMQTTClientOptions(pathconfig)

//Value1 - test
var Value1 = Test{
	Topic:   "/salt.coffee189@gmail.com/test",
	Qos:     0,
	Message: "Hello from the other side",
}

//Value2 - test
var Value2 = Test{
	Topic:   "/salt.coffee189@gmail.com/test",
	Qos:     1,
	Message: "It's over 9000!",
}

//Value3 - test
var Value3 = Test{
	Topic:   "/salt.coffee189@gmail.com/test",
	Qos:     2,
	Message: `{"car":"BWM i8"}`,
}
