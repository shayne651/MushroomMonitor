package main

import (
	"database/sql"
	"log"
	"net/http"

	mushroom_handler "github.com/shayne651/MushroomMonitor/internal/handler/mushroom"
	config_service "github.com/shayne651/MushroomMonitor/internal/services/config"
	message_service "github.com/shayne651/MushroomMonitor/internal/services/message"
	mushroom_service "github.com/shayne651/MushroomMonitor/internal/services/mushroom"

	stage_handler "github.com/shayne651/MushroomMonitor/internal/handler/stage"
	stage_service "github.com/shayne651/MushroomMonitor/internal/services/stage"

	grow_handler "github.com/shayne651/MushroomMonitor/internal/handler/grow"
	grow_service "github.com/shayne651/MushroomMonitor/internal/services/grow"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := config_service.Config{}
	config.Initialize()

	db := initializeDB()
	mq := initializeMQTT(config)
	mq.InitializeConnection("wkjnkjn")

	initializeRestAPI(db)
	select {}
}

func initializeDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./mush.db")
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

	//TODO: Add migrations and db maintenance to prune records over x days old
	dbMigrations(db)
	return db
}

func dbMigrations(db *sql.DB) {
	log.Println("Running DB Migrations")
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("Error getting db driver for migrations", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/",
		"sqlite3", driver)
	if err != nil {
		log.Fatal("Error migrating DB", err)
	}

	err = m.Up()
	if err != nil {
		log.Println("Error updating db")
		err = m.Down()
		if err != nil {
			log.Println("Error downgrading db")
		}
	}
	log.Println("Migration Successful")
}

func initializeMQTT(config config_service.Config) *message_service.MqttService {
	mq := message_service.MqttService{}
	mq.SetLoginDetails(config.MQTTHost, config.MQTTPort, config.MQTTUser, config.MQTTPass)
	return &mq
}

func initializeRestAPI(db *sql.DB) {
	mux := http.NewServeMux()
	ms := mushroom_service.MushroomService{DB: db}
	mh := mushroom_handler.MushroomHandler{MushroomService: &ms}
	mh.InitializeRestAPI(mux)

	initializeStage(db, mux)

	initializeGrow(db, mux)

	log.Panic(http.ListenAndServe(":7891", mux))
}

func initializeMushroom(db *sql.DB, mux *http.ServeMux) {

}

func initializeStage(db *sql.DB, mux *http.ServeMux) {
	stageService := stage_service.StageService{DB: db}
	stageHandler := stage_handler.StageHandler{StageService: &stageService}
	stageHandler.Initialize(mux)
}

func initializeGrow(db *sql.DB, mux *http.ServeMux) {
	growService := grow_service.GrowService{DB: db}
	growHandler := grow_handler.GrowHandler{GrowService: growService}
	growHandler.Initialize(mux)
}
