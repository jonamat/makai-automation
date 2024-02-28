package lights

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonamat/makai-automations/pkg/utils"
)

const ENABLED_TOPIC = "aut/light/enabled"
const LIGHT_LEVEL_TOPIC = "aut/light/level"

const DEFAULT_ENABLED = true
const DEFAULT_LIGHT_LEVEL = 255

var (
	mqttClient mqtt.Client
	dbClient   *badger.DB
)

func StartService() {
	fmt.Println("Starting lights job...")

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
	// getters and setters
	mqttClient.Subscribe(LIGHT_LEVEL_TOPIC+"/set", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		level, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			fmt.Println("Error parsing light level: ", err)
			return
		}

		err = setLightLevel(level)
		if err != nil {
			fmt.Println("Error setting light level: ", err)
			return
		}

		mqttClient.Publish(LIGHT_LEVEL_TOPIC, 0, false, fmt.Sprintf("%d", level))
	})

	// automation tasks
	mqttClient.Subscribe("dev/cabin-pir", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			level, err := getLightLevel()
			if err != nil {
				fmt.Println("Error getting light level: ", err)
				return
			}
			mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", level))
		case "OFF":
			mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", 0))
		}
	})
}

func disable(mqttClient mqtt.Client) {
	// termination tasks
	mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", 0))

	// unsubscribe to automation tasks
	mqttClient.Unsubscribe("dev/cabin-pir")
}
