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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DispatcherFunc func(chan bool) error

type DispatcherCommand struct {
	fn   DispatcherFunc
	done chan(bool)
}

var commands []DispatcherCommand
var done chan bool = make(chan bool)

// dispatchCmd represents the dispatch command
var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch services based on the config file.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("dispatch starting")
		defer log.Println("dispatch exiting on signal (eg. Ctrl-C)")
		srvs := viper.GetStringSlice("dispatch.services")
		log.Println("services:", srvs)
		for s := range srvs {
			switch srvs[s] {
			case "smstools":
				spawn(DispatcherFunc(smstools))
			}
		}
		log.Println("dispatch done - Ctrl-C to exit")
		<- done
		done <- true
	},
}

func spawn(fn DispatcherFunc) {
	dc := DispatcherCommand{fn: fn}
	dc.done = make(chan bool)
	go dc.fn(dc.done)
	commands = append(commands, dc)
}

func init() {
	RootCmd.AddCommand(dispatchCmd)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func(){
		<-c
		for c := range commands {
			commands[c].done <- true
			<- commands[c].done
		}
		done <- true
		<- done
		os.Exit(1)
	}()

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dispatchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dispatchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
