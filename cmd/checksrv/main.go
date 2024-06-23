package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

func main() {
	hasErr := false

	user, ok := os.LookupEnv("IOT_USER")
	if !ok {
		fmt.Println("Error: Environment variable IOT_USER not set")
		hasErr = true
	}

	pass, ok := os.LookupEnv("IOT_PASS")
	if !ok {
		fmt.Println("Error: Environment variable IOT_PASS not set")
		hasErr = true
	}

	addr, ok := os.LookupEnv("MQTT_ADDRESS")
	if !ok {
		fmt.Println("Error: Environment variable MQTT_ADDRESS not set")
		hasErr = true
	}

	if hasErr {
		os.Exit(1)
	}

	scheme, ok := os.LookupEnv("MQTT_SCHEME")
	if !ok {
		scheme = "tcp"
	}

	port, ok := os.LookupEnv("MQTT_PORT")
	if !ok {
		port = "1883"
	}

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("%s://%s:%s", scheme, addr, port))
	opts.SetUsername(user).SetPassword(pass)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("connected")
}
