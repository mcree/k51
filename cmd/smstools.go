// Copyright © 2016 Erno Rigo <erno@rigo.info>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/mcree/k51/backend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
	"log"
	"os"
	"github.com/eclipse/paho.mqtt.golang"
)

// smsCmd represents the sms command
var smsCmd = &cobra.Command{
	Use:   "smstools",
	Short: "Queue management for smstools",
	Long:  `Connects incoming and outgoing smstools3 daemon queues to MQTT`,
	Run: func(cmd *cobra.Command, args []string) {
		smstools(make(chan bool)) // no exit
	},
}

func init() {
	RootCmd.AddCommand(smsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// smsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// smsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func smstools(done chan bool) (error) {
	log.Println("smstools queue management starting")
	defer log.Println("smstools queue management exiting")

	var err error

	inChannel := viper.GetString("mqtt.channel") + "/sms/in"
	outChannel := viper.GetString("mqtt.channel") + "/sms/out"

	mq := backend.MQClient()
	mq.Publish(viper.GetString("mqtt.channel") + "/sms", 0, false, "test message").WaitTimeout(time.Second * 2)

	outd, err := backend.NewQueueDirWriter(viper.GetString("smstools.outgoing"),"sms_","")
	if err != nil {
		return err
	}
	defer outd.Close()

	ind, err := backend.NewQueueDirReader(viper.GetString("smstools.incoming"), func(c backend.QueueItem) {
		log.Println("Incoming sms: "+c.Name)
		mq.Publish(inChannel, 2, false, c.Payload)
		os.Remove(c.Name)
	})
	if err != nil {
		return err
	}
	defer ind.Close()

	mq.Subscribe(outChannel, 0, func(client mqtt.Client, msg mqtt.Message) {
		name, _ := outd.Write(msg.Payload())
		log.Println("Outgoing sms: "+name)
	} )
	defer mq.Unsubscribe(outChannel)

	<- done
	done <- true

	return err
}
