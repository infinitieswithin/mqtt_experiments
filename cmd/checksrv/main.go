package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"time"
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
	opts.SetDefaultPublishHandler(pubHandler)
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("connected")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("connection lost: %v\n", err)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	breakChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{}, 1)

	signal.Notify(breakChan, os.Interrupt, os.Kill)

	go pubWorker(client, doneChan)

	client.Subscribe("iot/test", 0, func(c mqtt.Client, m mqtt.Message) {
		fmt.Printf("[sub] received message on topic '%s': %s\n", m.Topic(), string(m.Payload()))
	})

	for sig := range breakChan {
		fmt.Println(sig, "received")
		doneChan <- struct{}{}
		client.Disconnect(250)
		break
	}

	fmt.Println("Done.")
}

func pubHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func pubWorker(client mqtt.Client, exitChan chan struct{}) {
	delayChan := time.After(time.Second)
	for {
		select {
		case <-exitChan:
			fmt.Println("terminating publish worker")
			return
		case <-delayChan:
			if token := client.Publish("iot/test", 2, false, "this is a test message"); token.Wait() && token.Error() != nil {
				fmt.Println("error publishing message", token.Error())
			}
			fmt.Println("[pub] message sent")
			delayChan = time.After(time.Second)
		}
	}
}
