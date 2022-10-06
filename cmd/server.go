/*
Copyright © 2022 clemon

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/c1emon/barkbridge/bridge"
	"github.com/c1emon/barkbridge/providers"
	"github.com/spf13/cobra"
)

var serverAddress string
var conf providers.AmqpConf

// var daemon bool

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start brak bridge server",
	Run: func(cmd *cobra.Command, args []string) {
		amqpProvider := providers.NewAmqpProvider(&conf)

		b := bridge.New(serverAddress)
		b.AddProvider("amqpProvider", amqpProvider)
		b.Server()

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&serverAddress, "bark-server", "b", "http://127.0.0.1:8080", "Bark server address")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// serverCmd.PersistentFlags().BoolVar(nil, "amqp", false, "enable amqp provider")
	serverCmd.PersistentFlags().StringVar(&conf.Addr, "amqp.address", "amqp://user:pass@127.0.0.1:5672", "amqp server address")
	serverCmd.PersistentFlags().StringVar(&conf.Exchange, "amqp.exchange", "amq.topic", "amqp exchange name")
	serverCmd.PersistentFlags().StringVar(&conf.Queue, "amqp.queue", "bark-bridge", "amqp queue name")
	serverCmd.PersistentFlags().StringVar(&conf.RoutingKey, "amqp.routingkey", "iot.sms.upload", "amqp routing key")
	serverCmd.MarkFlagsRequiredTogether("amqp.address", "amqp.exchange", "amqp.queue", "amqp.routingkey")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
