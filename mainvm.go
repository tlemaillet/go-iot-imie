package main

import (
	"fmt"
	"os"
	"os/signal"
	"encoding/json"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"github.com/franela/goreq"
	"./common"
)

func main() {
	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address: common.IpMosquitoServ + ":1883",
		ClientID: []byte("vm-client"),
	})
	if err != nil {
		panic(err)
	}

	// Subscribe to topics.
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("temp"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: sensorDataHandler,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Wait for receiving a signal.
	<-sigc

	// Unsubscribe from topics.
	err = cli.Unsubscribe(&client.UnsubscribeOptions{
		TopicFilters: [][]byte{
			[]byte("temp"),
		},
	})
	if err != nil {
		panic(err)
	}

	// Disconnect the Network Connection.
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
func sensorDataHandler(topicName, message []byte) {

	fmt.Println(string(topicName), string(message))
	if (string(topicName) != "temp") {
		panic("pas le bon message")
	}

	data := &common.SensorData{}
	err := json.Unmarshal(message, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonitem, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	request := goreq.Request{
		Method: "POST",
		Uri: "http://" + common.IpApiServ + "/cloud_api/web/app_dev.php/api/datasensors?sender=go",
		Accept: "application/json",
		ContentType: "application/json",
		UserAgent: "goreq",
		Body: string(jsonitem),
	}
	res, err := request.Do()
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Header)
	fmt.Println(res.Body.ToString())
}
