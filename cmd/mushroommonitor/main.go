package main

import (
	"database/sql"
	"log"

	config_service "github.com/shayne651/MushroomMonitor/interal/services/config"
	message_service "github.com/shayne651/MushroomMonitor/interal/services/message"

	_ "modernc.org/sqlite"
)

func main() {
	config := config_service.Config{}
	config.Initialize()

	_ = initializeDB()
	mq := initializeMQTT(config)
	mq.InitializeConnection("wkjnkjn")
}

func initializeDB() *sql.DB {
	db, err := sql.Open("sqlite", "../../mush.db")
	if err != nil {
		log.Panic("Connection to sqlite failed\n", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Panic("Connection started but cannot ping Sqlite\n", err)
		return nil
	}

	log.Println("DB connection successful")
	return db
}

func initializeMQTT(config config_service.Config) *message_service.MqttService {
	mq := message_service.MqttService{}
	mq.SetLoginDetails(config.MQTTHost, config.MQTTPort, config.MQTTUser, config.MQTTPass)
	return &mq
}
