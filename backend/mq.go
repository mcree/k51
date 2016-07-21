package backend

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"sync"
	"time"
	log "github.com/Sirupsen/logrus"
)

var client mqtt.Client
var onceMQ sync.Once

func defaultHandler(_ mqtt.Client, msg mqtt.Message) {
	log := log.WithFields(log.Fields{
		"prefix": "mq",
		"message": msg,
	})
	log.Error("default handler called")
}

func onConnectHandler(_ mqtt.Client) {
	log := log.WithField("prefix", "mq")
	log.Info("connected")
}

func onConnectionLostHandler(_ mqtt.Client, err error) {
	log := log.WithField("prefix", "mq")
	log.Error("connection lost: ", err)
}

// MQClient returns an initialized and connected, singleton MQTT client instance.
func MQClient() mqtt.Client {
	log := log.WithField("prefix", "mq")

	onceMQ.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(viper.GetString("mqtt.url"))
		opts.SetClientID(viper.GetString("mqtt.clientid"))
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(time.Second * 10)
		opts.SetDefaultPublishHandler(defaultHandler)
		opts.SetOnConnectHandler(onConnectHandler)
		opts.SetConnectionLostHandler(onConnectionLostHandler)
		opts.SetProtocolVersion(3)
		//fs := mqtt.NewFileStore(viper.GetString("mqtt.store"))
		//fs.Open()
		//opts.SetStore(fs)
		opts.SetStore(mqtt.NewFileStore(viper.GetString("mqtt.store")))
		log.Info("connecting: ", viper.GetString("mqtt.url"))
		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Panic(token.Error())
		}
	})
	return client
}

func MQCleanup() {
	log := log.WithField("prefix", "mq")
	MQClient().Disconnect(0)
	log.Info("disconnected")
}