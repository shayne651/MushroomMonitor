package message_service

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IMqttService interface {
	InitializeConnection(clientID string)
	PublishMessage(topic string, message []byte)
	SubscribeToTopic(topic string, handler func(client mqtt.Client, msg mqtt.Message))
}

type MqttService struct {
	host, port, username, password string
}

func (ms *MqttService) SetLoginDetails(host, port, username, password string) {
	ms.host = host
	ms.port = port
	ms.username = username
	ms.password = password
}

func (ms *MqttService) InitializeConnection(clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + ms.host + ":" + ms.port)
	opts.SetClientID(clientID)
	opts.SetUsername(ms.username)
	opts.SetPassword(ms.password)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Print("Error connecting to mosquitto\n", token.Error())
		return nil
	}
	return client
}

func (ms *MqttService) PublishMessage(topic string, message []byte) {
	client := ms.InitializeConnection("main-pub-" + topic)
	token := client.Publish(topic, 0, false, message)
	go func() {
		token.Wait()
		client.Disconnect(250)
	}()
}

func (ms *MqttService) SubscribeToTopic(topic string, handler func(client mqtt.Client, msg mqtt.Message)) {
	client := ms.InitializeConnection("main-reading-" + topic)

	go func() {
		token := client.Subscribe(topic, 0, handler)

		token.Wait()
	}()

}
