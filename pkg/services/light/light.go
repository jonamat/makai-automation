package lights

import (
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonamat/makai-automations/pkg/utils"
)

const AUTOMATION_TOPIC = "aut/light"
const DEFAULT_ENABLED = true
const DEFAULT_LIGHT_LEVEL = 255

var (
	mqttClient mqtt.Client
	dbClient   *badger.DB
)

func StartService() {
	fmt.Println("Starting lights job...")

	mqttClient = utils.CreateMqttClient()
	dbClient = utils.CreateDbClient()

	mqttClient.Subscribe(AUTOMATION_TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			setEnabled(true)
			enable(mqttClient)
		case "OFF":
			setEnabled(false)
			disable(mqttClient)
		}
	})

	// first start
	if getEnabled() {
		enable(mqttClient)
	}

	fmt.Println("Listening for messages on dev/cabin-pir")
	select {}
}

func enable(mqttClient mqtt.Client) {
	// getters and setters
	mqttClient.Subscribe("aut/light/level", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		level, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			fmt.Println("Error parsing light level: ", err)
			return
		}

		setLightLevel(level)
	})

	// automation tasks
	mqttClient.Subscribe("dev/cabin-pir", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", getLightLevel()))
		case "OFF":
			mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", 0))
		}
	})
}

func disable(mqttClient mqtt.Client) {
	// termination tasks
	mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", 0))

	// unsubscribtions
	mqttClient.Unsubscribe("dev/cabin-pir")
	mqttClient.Unsubscribe("aut/light/level")
}
