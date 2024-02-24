package lights

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonamat/makai-automations/pkg/utils"
)

func StartService() {
	fmt.Println("Starting lights job...")

	mqttClient := utils.CreateMqttClient()

	mqttClient.Subscribe("dev/cabin-pir", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		if string(msg.Payload()) == "ON" {
			mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", 255))

		} else {
			mqttClient.Publish("dev/main-light/set", 0, false, fmt.Sprintf("SET %d", 0))
		}
	})

	fmt.Println("Listening for messages on dev/cabin-pir")
	select {}
}
