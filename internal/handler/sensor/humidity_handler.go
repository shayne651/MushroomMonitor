package sensor_handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	hh.MQ.SubscribeToTopic("humidity-confirm", hh.handleConfirmHumidity)

	mux.HandleFunc("GET /humidity/{id}/{days}", hh.getAllHumidity)
	mux.HandleFunc("GET /humidity/last/{id}", hh.getLastHumidity)
}

func (hh *HumidityHandler) getAllHumidity(w http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	days, err := strconv.Atoi(request.PathValue("days"))

	if err != nil {
		log.Println("Error parsing days", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid days"))
		return
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing uuid", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid uuid"))
		return
	}

	humidity, err := hh.Hs.GetRecentHumidity(days, idParsed)
	humidityJson, err := json.Marshal(humidity)
	if err != nil {
		log.Println("Error marshaling humidity", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(humidityJson)
}

func (hh *HumidityHandler) getLastHumidity(w http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	idParsed, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing uuid", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid uuid"))
		return
	}
	humidity, err := hh.Hs.GetLastHumidity(idParsed)
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
	body := msg.Payload()
	err := json.Unmarshal(body, &humidity)
	if err != nil {
		log.Println("Error getting humidity from mqtt", err)
	}
	humidity.HumidityID = uuid.New()
	hh.Hs.SaveHumidity(humidity)
}
