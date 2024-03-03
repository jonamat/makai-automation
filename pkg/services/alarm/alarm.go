package alarm

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonamat/makai-automations/pkg/utils"
)

const ENABLED_TOPIC = "aut/alarm/enabled"

const DEFAULT_ENABLED = false
const DEFAULT_REGRET_TIME = 5 * time.Second

var (
	mqttClient     mqtt.Client
	dbClient       *badger.DB
	alarmIsCycling = false
)

var stopCycle = make(chan bool)

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
		case "STATE":
			enabled, _ := getEnabled()
			state := "OFF"
			if enabled {
				state = "ON"
			}

			mqttClient.Publish(ENABLED_TOPIC, 0, false, state)
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
	mqttClient.Subscribe("dev/cabin-door-sensor", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			go alarmCycle(stopCycle)
		}
	})

	mqttClient.Subscribe("dev/van-doors", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic: ", msg.Topic(), " with payload: ", string(msg.Payload()))

		switch string(msg.Payload()) {
		case "ON":
			go alarmCycle(stopCycle)
		}
	})
}

func disable(mqttClient mqtt.Client) {
	// termination tasks
	if alarmIsCycling {
		stopCycle <- true
	}

	// turn off the alarm if it's on (just in case)
	unFire()

	// unsubscribe to automation tasks
	mqttClient.Unsubscribe("dev/cabin-door-sensor")
	mqttClient.Unsubscribe("dev/van-doors")
}

func alarmCycle(stopCycle chan bool) {
	if alarmIsCycling {
		return
	} else {
		alarmIsCycling = true
	}

	// First signal the alarm armed with a 1/2 second fire
	timer := time.NewTimer(500 * time.Millisecond)
	fire()
	select {
	case <-timer.C:
		break
	case <-stopCycle:
		alarmIsCycling = false
		return
	}
	unFire()

	// Wait DEFAULT_REGRET_TIME seconds to regret the chooses of your life
	timer.Reset(DEFAULT_REGRET_TIME)
	select {
	case <-timer.C:
		break
	case <-stopCycle:
		alarmIsCycling = false
		return
	}

	// Fire on for 2 minutes
	fire()
	utils.SendSMS("Makai alarm fired!")
	timer.Reset(2 * time.Minute)
	select {
	case <-timer.C:
		break
	case <-stopCycle:
		alarmIsCycling = false
		unFire()
		return
	}

	unFire()
	alarmIsCycling = false
}

func fire() {
	fmt.Println("🔔 Alarm fired!")
	mqttClient.Publish("dev/cabin-alarm/set", 0, false, "ON")
	mqttClient.Publish("dev/van-alarm/set", 0, false, "ON")

}

func unFire() {
	fmt.Println("🔕 Alarm unfired!")
	mqttClient.Publish("dev/cabin-alarm/set", 0, false, "OFF")
	mqttClient.Publish("dev/van-alarm/set", 0, false, "OFF")

}
