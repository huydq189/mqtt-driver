package transport

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	errStillConnected   = errors.New("mqtt: still connected. Kill all processes manually")
	errMQTTDisconnected = errors.New("mqtt: disconnected")
)

// MQTT is the implementation of the MQTT client
type MQTT struct {
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
	connecting      bool
	disconnected    bool
	client          mqtt.Client
	opts            *mqtt.ClientOptions
	connectToken    *mqtt.ConnectToken
}

func (m *MQTT) getProtocol() string {
	if m.sslEnabled == true {
		return "ssl"
	}
	return "tcp"
}

//InitMQTTClientOptions - init MQTT options
func InitMQTTClientOptions() *MQTT {
	viper.SetConfigFile(`./config/config.json`)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	//DEBUG - Debugging
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//ERROR - Debugging
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	connID := uuid.New().String()

	mqttClient := &MQTT{
		host:            viper.GetString(`mqtt.host`),
		port:            viper.GetInt(`mqtt.port`),
		prefix:          viper.GetString(`mqtt.prefix`),
		clientID:        viper.GetString(`mqtt.clientID`) + "-" + connID,
		subscriptionQos: byte(viper.GetInt(`mqtt.subscriptionQos`)),
		persistent:      viper.GetBool(`mqtt.persistent`),
		order:           viper.GetBool(`mqtt.order`),
		sslEnabled:      viper.GetBool(`mqtt.sslEnabled`),
		username:        viper.GetString(`mqtt.username`),
		password:        viper.GetString(`mqtt.password`),
		clientCertPath:  viper.GetString(`mqtt.clientCertPath`),
		privateKeyPath:  viper.GetString(`mqtt.privateKeyPath`),
		keepAliveSec:    viper.GetInt(`mqtt.keepAliveSec`),
		pingTimeoutSec:  viper.GetInt(`mqtt.pingTimeoutSec`),
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
	mqttClient.opts.SetCleanSession(!mqttClient.persistent)
	mqttClient.opts.SetOrderMatters(mqttClient.order)
	mqttClient.opts.SetKeepAlive(time.Duration(mqttClient.keepAliveSec) * time.Second)
	mqttClient.opts.SetPingTimeout(time.Duration(mqttClient.pingTimeoutSec) * time.Second)
	mqttClient.opts.SetConnectionLostHandler(mqttClient.connectionLostHandler)
	return mqttClient
}

//CreateMQTTClient - Create a new MQTT client
func (m *MQTT) CreateMQTTClient() *MQTT {
	if m.opts == nil {
		log.Fatalln("No options created")
		return nil
	}
	m.client = mqtt.NewClient(m.opts)
	return m
}

// Start running the MQTT client
func (m *MQTT) Start() {
	log.Printf("Starting MQTT client on %s://%s:%v with Persistence:%v, OrderMatters:%v, KeepAlive:%v, PingTimeout:%v, QOS:%v\n",
		m.getProtocol(), m.host, m.port, m.persistent, m.order, m.keepAliveSec, m.pingTimeoutSec, m.subscriptionQos)
	m.connect()
}

func (m *MQTT) connect() {
	m.connectToken = m.client.Connect().(*mqtt.ConnectToken)
	if m.connectToken.Wait() && m.connectToken.Error() != nil {
		if !m.connecting {
			log.Printf("MQTT client %s", m.connectToken.Error())
			m.retryConnect()
		}
	} else {
		fmt.Println("Connected with broker")
	}
}

// Stop the MQTT client
func (m *MQTT) Stop() error {
	m.client.Disconnect(500)
	if m.client.IsConnected() {
		return errStillConnected
	}
	return nil
}

// Subscribe - subcribe to a topic
func (m *MQTT) Subscribe(topic string, handler mqtt.MessageHandler, qos byte) error {
	token := m.client.Subscribe(topic, qos, handler)
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// SubscribeOC - Set options subcribe to a topic onconnect
func (m *MQTT) SubscribeOC(topic string, handler mqtt.MessageHandler, qos byte) {
	var f = func(c mqtt.Client) {
		token := m.client.Subscribe(topic, m.subscriptionQos, handler)
		if token.WaitTimeout(2) && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}
	m.opts.SetOnConnectHandler(f)
}

// SubscribeD - SubscribeDefault to a topic with QOS
func (m *MQTT) SubscribeD(topic string) error {
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
func (m *MQTT) UnSubscribe(topic string) error {
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
func (m *MQTT) Publish(topic, message string) error {
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	var q = m.subscriptionQos
	// Send - Publish a message to Topic
	token := m.client.Publish(topic, q, false, message)
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// PublishQOS - qos
func (m *MQTT) PublishQOS(topic, message string, qos byte) error {
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	// Send - Publish a message to Topic
	token := m.client.Publish(topic, qos, false, message)
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// RetryConnect - reconnect after time
func (m *MQTT) retryConnect() {
	log.Printf("MQTT client starting reconnect procedure in background")
	m.connecting = true
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			m.connect()
			if m.client.IsConnected() {
				ticker.Stop()
				m.connecting = false
			}
		}
	}()
}

//call retryConnect
func (m *MQTT) connectionLostHandler(c mqtt.Client, err error) {
	log.Printf("MQTT client lost connection: %v", err)
	m.disconnected = true
	m.retryConnect()
}
