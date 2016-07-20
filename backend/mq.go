package backend

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var client mqtt.Client
var onceMQ sync.Once

// MQClient returns an initialized and connected, singleton MQTT client instance.
func MQClient() mqtt.Client {
	onceMQ.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(viper.GetString("mqtt.url"))
		opts.SetClientID(viper.GetString("mqtt.clientid"))
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(time.Second * 10)
		//fs := mqtt.NewFileStore(viper.GetString("mqtt.store"))
		//fs.Open()
		//opts.SetStore(fs)
		//opts.SetStore(mqtt.NewFileStore(viper.GetString("mqtt.store")))
		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	})
	return client
}
