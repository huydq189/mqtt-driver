package core

import (
	"errors"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	errStillConnected   = errors.New("mqtt: still connected. Kill all processes manually")
	errMQTTDisconnected = errors.New("mqtt: disconnected")
)

// MQTT - IModelRepository
type MQTT interface {
	Start() (bool, error)
	Stop() error
	Option() *mqtt.ClientOptions
	Ping() bool
	SetOnConnectHandler(handler mqtt.OnConnectHandler) *mqtt.ClientOptions
	Subscribe(topic string, handler mqtt.MessageHandler, qos byte) error
	//SubscribeOC(topic string, handler mqtt.MessageHandler, qos byte) *mqtt.ClientOptions
	SubscribeD(topic string) error
	UnSubscribe(topic string) error
	Publish(topic string, payload interface{}, retained bool) error
	PublishQOS(topic string, payload interface{}, retained bool, qos byte) error
}

// MQTTClient is the implementation of the MQTT client
type mqttClient struct {
	host            string
	port            int
	prefix          string
	clientID        string
	sslEnabled      bool
	username        string
	password        string
	clientCertPath  string
	privateKeyPath  string
	keepAliveSec    int
	pingTimeoutSec  int
	subscriptionQos byte
	persistent      bool
	order           bool
	willEnabled     bool
	willTopic       string
	willPayload     string
	willQos         byte
	willRetained    bool
	connecting      bool
	disconnected    bool
	resumesub       bool
	timeout         int
	client          mqtt.Client
	opts            *mqtt.ClientOptions
	connectToken    *mqtt.ConnectToken
}

func (m *mqttClient) getProtocol() string {
	if m.sslEnabled {
		return "ssl"
	}
	return "tcp"
}

//InitMQTTClientOptions - Init MQTT options
func InitMQTTClientOptions(path string) MQTT {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	// DEBUG - Debugging
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	// ERROR - Debugging
	// mqtt.ERROR = log.New(os.Stdout, "", 0)
	connID := uuid.New().String()
	mqttClient := &mqttClient{
		host:            viper.GetString(`mqtt.host`),
		port:            viper.GetInt(`mqtt.port`),
		prefix:          viper.GetString(`mqtt.prefix`),
		clientID:        viper.GetString(`mqtt.clientID`) + "-" + connID,
		subscriptionQos: byte(viper.GetInt(`mqtt.subscriptionQos`)),
		persistent:      viper.GetBool(`mqtt.persistent`),
		order:           viper.GetBool(`mqtt.order`),
		willEnabled:     viper.GetBool(`mqtt.WillEnabled`),
		willTopic:       viper.GetString(`mqtt.WillTopic`),
		willPayload:     viper.GetString(`mqtt.WillPayload`),
		willQos:         byte(viper.GetInt(`mqtt.WillQos`)),
		willRetained:    viper.GetBool(`mqtt.WillRetained`),
		sslEnabled:      viper.GetBool(`mqtt.sslEnabled`),
		username:        viper.GetString(`mqtt.username`),
		password:        viper.GetString(`mqtt.password`),
		clientCertPath:  viper.GetString(`mqtt.clientCertPath`),
		privateKeyPath:  viper.GetString(`mqtt.privateKeyPath`),
		keepAliveSec:    viper.GetInt(`mqtt.keepAliveSec`),
		pingTimeoutSec:  viper.GetInt(`mqtt.pingTimeoutSec`),
		timeout:         viper.GetInt(`mqtt.timeout`),
		resumesub:       viper.GetBool(`mqtt.resumesub`),
	}
	mqttClient.opts = mqtt.NewClientOptions()
	// CreateMQTTClient creates a new MQTT client
	if mqttClient.username != "" {
		mqttClient.opts.SetUsername(mqttClient.username)
	}
	if mqttClient.password != "" {
		mqttClient.opts.SetPassword(mqttClient.password)
	}
	mqttClient.opts.AddBroker(fmt.Sprintf("%s://%s:%v", mqttClient.getProtocol(), mqttClient.host, mqttClient.port))
	mqttClient.opts.SetClientID(mqttClient.clientID)
	//mqttClient.opts.SetAutoReconnect(true)
	mqttClient.opts.SetCleanSession(!mqttClient.persistent)
	mqttClient.opts.SetOrderMatters(mqttClient.order)
	if mqttClient.willEnabled {
		mqttClient.opts.SetWill(mqttClient.willTopic, mqttClient.willPayload, mqttClient.willQos, mqttClient.willRetained)
	}
	mqttClient.opts.SetAutoReconnect(false)
	mqttClient.opts.SetKeepAlive(time.Duration(mqttClient.keepAliveSec) * time.Second)
	mqttClient.opts.SetPingTimeout(time.Duration(mqttClient.pingTimeoutSec) * time.Second)
	mqttClient.opts.SetResumeSubs(mqttClient.resumesub)
	mqttClient.opts.SetConnectionLostHandler(mqttClient.connectionLostHandler)
	return mqttClient
}

//SetOption - Set the option before start
func (m *mqttClient) Option() *mqtt.ClientOptions {
	return m.opts
}

//Ping - check is connected
func (m *mqttClient) Ping() bool {
	return m.client.IsConnected()
}

//SetOnConnectHandler - Set handler to resume your action after lost connection if persistent is false
func (m *mqttClient) SetOnConnectHandler(handler mqtt.OnConnectHandler) *mqtt.ClientOptions {
	m.opts = m.opts.SetOnConnectHandler(handler)
	return m.opts
}

// Start running the MQTT client
func (m *mqttClient) Start() (bool, error) {
	//CreateMQTTClient - Create a new MQTT client
	if m.opts == nil {
		return false, errors.New("No opstion defined")
	}
	if m.opts.OnConnect == nil && m.opts.CleanSession {
		log.Println("OnConnect hasn't set yet: Please set OnConnectHandler if you want to resume your action ex:subcription(subscribe to topic) after lost connection ")
	}
	m.client = mqtt.NewClient(m.opts)
	log.Printf("Starting MQTT client %s on %s://%s:%v with Persistence:%v, OrderMatters:%v, KeepAlive:%v, PingTimeout:%v, QOS:%v, Resumesubs:%v\n", m.opts.ClientID,
		m.getProtocol(), m.host, m.port, !m.opts.CleanSession, m.opts.Order, m.opts.KeepAlive, m.opts.PingTimeout, m.subscriptionQos, m.opts.ResumeSubs)
	boo, err := m.connect()
	return boo, err
}

func (m *mqttClient) connect() (bool, error) {
	m.connectToken = m.client.Connect().(*mqtt.ConnectToken)
	var res bool
	var err error
	if m.connectToken.Wait() && m.connectToken.Error() != nil {
		if !m.connecting {
			log.Printf("MQTT client %s", m.connectToken.Error())
			res, err = m.retryConnect()

		}
	} else {
		fmt.Println("Connected with broker")
		res, err = true, nil
	}
	return res, err
}

// Stop the MQTT client
func (m *mqttClient) Stop() error {
	m.client.Disconnect(500)
	if m.client.IsConnected() {
		return errStillConnected
	}
	return nil
}

// Subscribe - subcribe to a topic
func (m *mqttClient) Subscribe(topic string, handler mqtt.MessageHandler, qos byte) error {
	token := m.client.Subscribe(topic, qos, handler)
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// SubscribeOC - Set options subcribe to a topic onconnect
// func (m *MQTT) SubscribeOC(topic string, handler mqtt.MessageHandler, qos byte) *mqtt.ClientOptions {
// 	var f = func(c mqtt.Client) {
// 		token := m.client.Subscribe(topic, m.subscriptionQos, handler)
// 		if token.WaitTimeout(2) && token.Error() != nil {
// 			log.Fatal(token.Error())
// 		}
// 	}
// 	m.opts = m.opts.SetOnConnectHandler(f)
// 	return m.opts
// }

// SubscribeD - SubscribeDefault to a topic with QOS
func (m *mqttClient) SubscribeD(topic string) error {
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	token := m.client.Subscribe(topic, m.subscriptionQos, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("[%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// UnSubscribe - Unsubcribe topic
func (m *mqttClient) UnSubscribe(topic string) error {
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	token := m.client.Unsubscribe(topic)
	if token.WaitTimeout(3) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Publish - Send messages to topic
func (m *mqttClient) Publish(topic string, payload interface{}, retained bool) error {
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	var q = m.subscriptionQos
	// Send - Publish a message to Topic
	token := m.client.Publish(topic, q, retained, payload)
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// PublishQOS - qos
func (m *mqttClient) PublishQOS(topic string, payload interface{}, retained bool, qos byte) error {
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	// Send - Publish a message to Topic
	token := m.client.Publish(topic, qos, retained, payload)
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// RetryConnect - reconnect after time
func (m *mqttClient) retryConnect() (bool, error) {
	log.Printf("MQTT client starting reconnect procedure in background")
	m.connecting = true
	timeout := time.After(time.Duration(m.timeout) * time.Second)
	tick := time.Tick(5000 * time.Millisecond)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			log.Fatalf("Client tried to connect but time out and failed")
			return false, errors.New("timed out")
		// Got a tick, we should check on doSomething()
		case <-tick:
			m.connect()
			if m.client.IsConnected() {
				return true, nil
			}
			// doSomething() didn't work yet, but it didn't fail, so let's try again
			// this will exit up to the for loop
		}
	}
}

//call retryConnect
func (m *mqttClient) connectionLostHandler(c mqtt.Client, err error) {
	log.Printf("MQTT client lost connection: %v", err)
	m.disconnected = true
	m.retryConnect()
}
