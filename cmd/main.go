package main

import (
	"os"

	"github.com/joho/godotenv"
	alarm "github.com/jonamat/makai-automations/pkg/services/alarm"
	clock "github.com/jonamat/makai-automations/pkg/services/clock"
	energy "github.com/jonamat/makai-automations/pkg/services/energy"
	light "github.com/jonamat/makai-automations/pkg/services/light"
)

func main() {
	println("Starting automation services...")

	// load dotenv file, but do not override existing environment variables
	err := godotenv.Load()
	if err != nil {
		println("Error loading .env file")
	}

	// check if all necessary environment variables are set
	keys := []string{
		"MQTT_BROKER",
		"MQTT_PORT",
	}
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); !ok || value == "" {
			println("Error: Environment variable not set: ", key)
			os.Exit(1)
		}
	}

	// Start all services
	go alarm.StartService()
	go clock.StartService()
	go energy.StartService()
	go light.StartService()

	// keep the main process running
	select {}
}
