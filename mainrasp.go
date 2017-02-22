package main

import (
	"fmt"
	"os"
	//"os/exec"
	//"log"
	"os/signal"

	"encoding/json"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"github.com/franela/goreq"
	"./common"
)

 

func main() {
	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	//// send two pings and send the ouput to STDOUT
	//output1, err := exec.Command("ping", "-c", "2", "8.8.8.8").Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("The results are -->  %s\n", output1)
	//
	//// run the date command and store the output
	//output2, err := exec.Command("date").Output()
	//if err != nil {
	//	log.Print(err)
	//}
	//// another wat to print
	//log.Printf("The date is --> %s ", output2)
	sensor := &common.Sensor{
		DisplayName: "Rasp-Sensor",
		Vendor: "Raspberry Foundation",
		Product: "Pi",
		Version: 3,
	}

	jsonsensor, err := json.Marshal(sensor)
	if err != nil {
		fmt.Println(err)
		return
	}

	request := goreq.Request{
		Method: "POST",
		Uri: "http://" + common.IpApiServ + "/cloud_api/web/app_dev.php/api/sensors?sender=go",
		Accept: "application/json",
		ContentType: "application/json",
		UserAgent: "goreq",
		Body: string(jsonsensor),
	}
	resguid, err := request.Do()
	if err != nil {
		panic(err)
	}

	fmt.Println(resguid.Header)
	fmt.Println(resguid.Body.ToString())

	var guid interface{}
	resguid.Body.FromJsonTo(guid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(guid)

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
	err = cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  common.IpMosquitoServ + ":1883",
		ClientID: []byte("rasp-client"),
	})
	if err != nil {
		panic(err)
	}

	data := &common.SensorData{
		SensorName:  "truc",
		Measurement: "temp",
		Time:        time.Now().UnixNano(),
		Value:       "OVER 9000",
	}

	jsonitem, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Publish a message.
	err = cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte("temp"),
		Message:   []byte(jsonitem),
	})
	if err != nil {
		panic(err)
	}

	// Wait for receiving a signal.
	<-sigc

	// Disconnect the Network Connection.
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
