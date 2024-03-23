package config

import (
	"log"
	"os"
)

type Config struct {
	MQTTHost, MQTTPort, MQTTUser, MQTTPass, MigrationLocation string
}

func (c *Config) Initialize() {
	local := os.Getenv("LOCAL_DEV")
	if local == "true" {
		log.Println("Initializing local dev configuration")
		c.MQTTHost = "localhost"
		c.MQTTPort = "1883"
		c.MQTTUser = "reefmonitor"
		c.MQTTPass = "reefmonitor"
		c.MigrationLocation = "db/migrations"
	} else {
		log.Println("Initializing Configuration")
		c.MQTTHost = os.Getenv("MQTT_HOST")
		c.MQTTPort = os.Getenv("MQTT_PORT")
		c.MQTTUser = os.Getenv("MQTT_USER")
		c.MQTTPass = os.Getenv("MQTT_PASSWORD")
		c.MigrationLocation = os.Getenv("MIGRATION_LOCATION")
	}
}
