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

type IFeaHandler interface {
}

type FeaHandler struct {
	Fs sensor_service.FeaService
	MQ *message_service.MqttService
}

func (fh *FeaHandler) InitializeFea(mux *http.ServeMux) {
	fh.MQ.SubscribeToTopic("fea-confirm", fh.handleConfirmFea)

	mux.HandleFunc("GET /fea/{id}/{days}", fh.getAllFea)
	mux.HandleFunc("GET /fea/last/{id}", fh.getLastFea)
}

func (fh *FeaHandler) getAllFea(w http.ResponseWriter, request *http.Request) {
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

	fea, err := fh.Fs.GetRecentFea(days, idParsed)
	feaJson, err := json.Marshal(fea)
	if err != nil {
		log.Println("Error marshaling fea", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(feaJson)
}

func (fh *FeaHandler) getLastFea(w http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	idParsed, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing uuid", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid uuid"))
		return
	}
	fea, err := fh.Fs.GetLastFea(idParsed)
	feaJson, err := json.Marshal(fea)
	if err != nil {
		log.Println("Error marshaling fea", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(feaJson)
}

func (fh *FeaHandler) handleConfirmFea(client mqtt.Client, msg mqtt.Message) {
	fea := sensor_service.Fea{}
	body := msg.Payload()
	err := json.Unmarshal(body, &fea)
	if err != nil {
		log.Println("Error getting fea from mqtt", err)
	}
	fea.FeaID = uuid.New()
	fh.Fs.SaveFea(fea)
}
