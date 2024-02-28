package alarm

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonamat/makai-automations/pkg/utils"
)

const ENABLED_TOPIC = "aut/alarm/enabled"

const DEFAULT_ENABLED = true

var (
	mqttClient     mqtt.Client
	dbClient       *badger.DB
	alarmIsCycling = false
)

func StartService() {
	fmt.Println("Starting alarm job...")

	mqttClient = utils.CreateMqttClient()
	defer mqttClient.Disconnect(250)
	dbClient = utils.CreateDbClient()
	defer dbClient.Close()

	for !mqttClient.IsConnected() {
		fmt.Println("Waiting for mqtt connection...")
		time.Sleep(1 * time.Second)
	}

	setupDb(dbClient)

	mqttClient.Subscribe(ENABLED_TOPIC+"/set", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			setEnabled(true)
			enable()
			mqttClient.Publish(ENABLED_TOPIC, 0, false, "ON")
		case "OFF":
			setEnabled(false)
			disable(mqttClient)
			mqttClient.Publish(ENABLED_TOPIC, 0, false, "OFF")
		}
	})

	// first start
	enabled, _ := getEnabled()
	if enabled {
		enable()
	}

	select {}
}

func enable() {
	// automation tasks
	// mqttClient.Subscribe("dev/cabin-pir", 0, func(client mqtt.Client, msg mqtt.Message) {
	// 	fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

	// 	switch string(msg.Payload()) {
	// 	case "ON":
	// 		alarmCycle()
	// 	}
	// })

	mqttClient.Subscribe("dev/cabin-door-sensor", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			alarmCycle()
		}
	})

	mqttClient.Subscribe("dev/van-doors", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			alarmCycle()
		}
	})
}

func disable(mqttClient mqtt.Client) {
	// termination tasks
	unFire()

	// unsubscribtions
	mqttClient.Unsubscribe(ENABLED_TOPIC + "/set")
}

func alarmCycle() {
	if alarmIsCycling {
		return
	} else {
		alarmIsCycling = true
	}

	// Alarm on for 1 second
	fire()
	time.Sleep(1 * time.Second)
	unFire()

	// Wait for 7 seconds to regret
	time.Sleep(7 * time.Second)

	enabled, _ := getEnabled()
	if !enabled {
		alarmIsCycling = false
		return
	}

	// Alarm on for 2 minutes
	fire()
	time.Sleep(2 * time.Minute)
	unFire()

	alarmIsCycling = false
}

func fire() {
	mqttClient.Publish("dev/cabin-alarm/set", 0, false, "ON")
	mqttClient.Publish("dev/van-alarm/set", 0, false, "ON")
}

func unFire() {
	mqttClient.Publish("dev/cabin-alarm/set", 0, false, "OFF")
	mqttClient.Publish("dev/van-alarm/set", 0, false, "OFF")
}
