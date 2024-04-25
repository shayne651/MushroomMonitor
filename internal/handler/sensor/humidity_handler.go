package sensor_handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	message_service "github.com/shayne651/MushroomMonitor/internal/services/message"
	sensor_service "github.com/shayne651/MushroomMonitor/internal/services/sensor"
)

type IHumidityHandler interface {
}

type HumidityHandler struct {
	Hs sensor_service.HumidityService
	MQ *message_service.MqttService
}

func (hh *HumidityHandler) InitializeHumidity(mux *http.ServeMux) {
	hh.MQ.SubscribeToTopic("mushroom_monitor-test.humidity", hh.handleConfirmHumidity)

	mux.HandleFunc("GET /humidity/{name}/{days}", hh.getAllHumidity)
	mux.HandleFunc("GET /humidity/last/{name}", hh.getLastHumidity)
}

func (hh *HumidityHandler) getAllHumidity(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")
	days, err := strconv.Atoi(request.PathValue("days"))

	if err != nil {
		log.Println("Error parsing days", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid days"))
		return
	}

	humidity, err := hh.Hs.GetRecentHumidity(days, name)
	humidityJson, err := json.Marshal(humidity)
	if err != nil {
		log.Println("Error marshaling humidity", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(humidityJson)
}

func (hh *HumidityHandler) getLastHumidity(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")

	humidity, _ := hh.Hs.GetLastHumidity(name)
	humidityJson, err := json.Marshal(humidity)
	if err != nil {
		log.Println("Error marshaling humidity", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(humidityJson)
}

func (hh *HumidityHandler) handleConfirmHumidity(client mqtt.Client, msg mqtt.Message) {
	humidity := sensor_service.Humidity{}
	humString := string(msg.Payload())

	h, err := strconv.ParseFloat(humString, 32)
	if err != nil {
		log.Println("Error parsing humidity", err)
	}
	humidity.HumidityID = uuid.New()
	humidity.Value = float32(h)
	humidity.RecordDate = int(time.Now().Unix())
	humidity.GrowName = "test"
	hh.Hs.SaveHumidity(humidity)
}
