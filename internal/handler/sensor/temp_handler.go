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

type ITempHandler interface {
}

type TempHandler struct {
	Ts sensor_service.TempService
	MQ *message_service.MqttService
}

func (th *TempHandler) InitializeTemp(mux *http.ServeMux) {
	th.MQ.SubscribeToTopic("mushroom_monitor-test.temp", th.handleConfirmTemp)

	mux.HandleFunc("GET /temp/{id}/{days}", th.getAllTemp)
	mux.HandleFunc("GET /temp/last/{id}", th.getLastTemp)
}

func (th *TempHandler) getAllTemp(w http.ResponseWriter, request *http.Request) {
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

	temp, err := th.Ts.GetRecentTemp(days, idParsed)
	tempJson, err := json.Marshal(temp)
	if err != nil {
		log.Println("Error marshaling temp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(tempJson)
}

func (th *TempHandler) getLastTemp(w http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	idParsed, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing uuid", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid uuid"))
		return
	}
	temp, err := th.Ts.GetLastTemp(idParsed)
	tempJson, err := json.Marshal(temp)
	if err != nil {
		log.Println("Error marshaling temp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(tempJson)
}

func (th *TempHandler) handleConfirmTemp(client mqtt.Client, msg mqtt.Message) {
	log.Println("Received temp comfirm")
	temp := sensor_service.Temp{}
	body := msg.Payload()

	err := json.Unmarshal(body, &temp)
	if err != nil {
		log.Println("Error getting temp from mqtt", err)
	}
	temp.TempID = uuid.New()
	th.Ts.SaveTemp(temp)
}
