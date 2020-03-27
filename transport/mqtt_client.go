package transport

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

// Client interface
type Client = mqtt.Client

// Message interface
type Message = mqtt.Message

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
	connectToken    *mqtt.ConnectToken
}

func (m *MQTT) getProtocol() string {
	if m.sslEnabled == true {
		return "ssl"
	}
	return "tcp"
}

func initMQTTClientOptions(m *MQTT) (*mqtt.ClientOptions, error) {

	opts := mqtt.NewClientOptions()

	if m.username != "" {
		opts.SetUsername(m.username)
	}
	if m.password != "" {
		opts.SetPassword(m.password)
	}
	opts.AddBroker(fmt.Sprintf("%s://%s:%v", m.getProtocol(), m.host, m.port))
	opts.SetClientID(m.clientID)
	opts.SetCleanSession(!m.persistent)
	opts.SetOrderMatters(m.order)
	opts.SetKeepAlive(time.Duration(m.keepAliveSec) * time.Second)
	opts.SetPingTimeout(time.Duration(m.pingTimeoutSec) * time.Second)
	opts.SetAutoReconnect(false)
	opts.SetConnectionLostHandler(m.connectionLostHandler)
	return opts, nil
}

// CreateMQTTClient creates a new MQTT client
func CreateMQTTClient() *MQTT {
	viper.SetConfigFile(`./config/config.json`)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
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
	opts, err := initMQTTClientOptions(mqttClient)
	if err != nil {
		log.Fatalf("unable to configure MQTT client: %s", err)
	}
	Client := mqtt.NewClient(opts)
	mqttClient.client = Client

	return mqttClient
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
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	token := m.client.Subscribe(topic, qos, handler)
	if token.WaitTimeout(2) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// SubscribeDefault - SubscribeDefault to a topic with QOS
func (m *MQTT) SubscribeDefault(topic string) error {
	if !m.client.IsConnected() {
		return errMQTTDisconnected
	}
	token := m.client.Subscribe(topic, m.subscriptionQos, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
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
			m.client.Connect()
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
