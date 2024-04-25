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

	sensor_handlers "github.com/shayne651/MushroomMonitor/internal/handler/sensor"
	sensor_service "github.com/shayne651/MushroomMonitor/internal/services/sensor"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := config_service.Config{}
	config.Initialize()

	db := initializeDB(config)
	mq := initializeMQTT(config)
	mq.InitializeConnection("wkjnkjn")

	initializeRestAPI(db, mq)
	select {}
}

func initializeDB(config config_service.Config) *sql.DB {
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
	dbMigrations(db, config)
	return db
}

func dbMigrations(db *sql.DB, config config_service.Config) {
	log.Println("Running DB Migrations")
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("Error getting db driver for migrations", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+config.MigrationLocation,
		"sqlite3", driver)
	if err != nil {
		log.Fatal("Error migrating DB", err)
		return
	}

	err = m.Up()
	if err != migrate.ErrNoChange && err != nil {
		log.Println("Error updating db: ", err)
		err = m.Down()
		if err != nil {
			log.Println("Error downgrading db: ", err)
			return
		}
		return
	}
	log.Println("Migration Successful")
}

func initializeMQTT(config config_service.Config) *message_service.MqttService {
	mq := message_service.MqttService{}
	mq.SetLoginDetails(config.MQTTHost, config.MQTTPort, config.MQTTUser, config.MQTTPass)
	return &mq
}

func initializeRestAPI(db *sql.DB, mq *message_service.MqttService) {
	mux := http.NewServeMux()

	initializeMushroom(db, mux)

	initializeStage(db, mux)

	initializeGrow(db, mq, mux)

	initializeTemp(db, mq, mux)

	initializeHumidity(db, mq, mux)

	initializeFea(db, mq, mux)

	log.Panic(http.ListenAndServe(":7891", mux))
}

func initializeMushroom(db *sql.DB, mux *http.ServeMux) {
	ms := mushroom_service.MushroomService{DB: db}
	mh := mushroom_handler.MushroomHandler{MushroomService: &ms}
	mh.InitializeRestAPI(mux)
}

func initializeStage(db *sql.DB, mux *http.ServeMux) {
	stageService := stage_service.StageService{DB: db}
	stageHandler := stage_handler.StageHandler{StageService: &stageService}
	stageHandler.Initialize(mux)
}

func initializeGrow(db *sql.DB, mq *message_service.MqttService, mux *http.ServeMux) {
	growService := grow_service.GrowService{DB: db}
	growHandler := grow_handler.GrowHandler{GrowService: growService, MQ: mq}
	growHandler.Initialize(mux)
}

func initializeTemp(db *sql.DB, mq *message_service.MqttService, mux *http.ServeMux) {
	tempService := sensor_service.TempService{DB: db}
	tempHandler := sensor_handlers.TempHandler{Ts: tempService, MQ: mq}
	tempHandler.InitializeTemp(mux)
}

func initializeHumidity(db *sql.DB, mq *message_service.MqttService, mux *http.ServeMux) {
	humidityService := sensor_service.HumidityService{DB: db}
	humidityHandler := sensor_handlers.HumidityHandler{Hs: humidityService, MQ: mq}
	humidityHandler.InitializeHumidity(mux)
}

func initializeFea(db *sql.DB, mq *message_service.MqttService, mux *http.ServeMux) {
	feaService := sensor_service.FeaService{DB: db}
	feaHandler := sensor_handlers.FeaHandler{Fs: feaService, MQ: mq}
	feaHandler.InitializeFea(mux)
}

// func initializeGraph(db *sql.DB, mux *http.ServeMux) {
// 	graphService := graph_service.GrowService{DB: db}
// 	graphHandler := graph_handler.GrowHandler{GrowService: graphService}
// 	graphHandler.Initialize(mux)
// }
