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

type ITempHandler interface {
}

type TempHandler struct {
	Ts sensor_service.TempService
	MQ *message_service.MqttService
}

func (th *TempHandler) InitializeTemp(mux *http.ServeMux) {
	th.MQ.SubscribeToTopic("mushroom_monitor-test.temp", th.handleConfirmTemp)

	mux.HandleFunc("GET /temp/{name}/{days}", th.getAllTemp)
	mux.HandleFunc("GET /temp/last/{name}", th.getLastTemp)
}

func (th *TempHandler) getAllTemp(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")
	days, err := strconv.Atoi(request.PathValue("days"))

	if err != nil {
		log.Println("Error parsing days", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid days"))
		return
	}
	temp, err := th.Ts.GetRecentTemp(days, name)
	tempJson, err := json.Marshal(temp)
	if err != nil {
		log.Println("Error marshaling temp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(tempJson)
}

func (th *TempHandler) getLastTemp(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")
	temp, err := th.Ts.GetLastTemp(name)
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

	tempString := string(msg.Payload())

	t, err := strconv.ParseFloat(tempString, 32)
	if err != nil {
		log.Println("Error parsing temp", err)
	}
	temp.TempID = uuid.New()
	temp.Value = float32(t)
	temp.RecordDate = int(time.Now().Unix())
	temp.GrowName = "test"
	th.Ts.SaveTemp(temp)
}
