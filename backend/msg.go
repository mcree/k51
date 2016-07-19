package backend

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var client mqtt.Client
var once sync.Once

// MQClient returns an initialized and connected, singleton MQTT client instance.
func MQClient() mqtt.Client {
	once.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(viper.GetString("mqtt.url"))
		opts.SetClientID(viper.GetString("mqtt.clientid"))
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(time.Second * 1)
		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		mqtt.MessageHandler()
	})
	return client
}
