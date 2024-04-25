package sensor_handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	// fh.MQ.SubscribeToTopic("fea-confirm", fh.handleConfirmFea)

	mux.HandleFunc("GET /fea/{name}/{days}", fh.getAllFea)
	mux.HandleFunc("GET /fea/last/{name}", fh.getLastFea)
}

func (fh *FeaHandler) getAllFea(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")
	days, err := strconv.Atoi(request.PathValue("days"))

	if err != nil {
		log.Println("Error parsing days", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid days"))
		return
	}

	fea, err := fh.Fs.GetRecentFea(days, name)
	feaJson, err := json.Marshal(fea)
	if err != nil {
		log.Println("Error marshaling fea", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(feaJson)
}

func (fh *FeaHandler) getLastFea(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")

	fea, err := fh.Fs.GetLastFea(name)
	feaJson, err := json.Marshal(fea)
	if err != nil {
		log.Println("Error marshaling fea", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(feaJson)
}

// func (fh *FeaHandler) handleConfirmFea(client mqtt.Client, msg mqtt.Message) {
// 	fea := sensor_service.Fea{}
// 	tempString := string(msg.Payload())

// 	t, err := strconv.ParseFloat(tempString, 32)
// 	if err != nil {
// 		log.Println("Error parsing temp", err)
// 	}
// 	fea.TempID = uuid.New()
// 	fea.Value = float32(t)
// 	fea.RecordDate = int(time.Now().Unix())
// 	fea.GrowName = "test"
// 	fh.Fs.SaveFea(fea)
// }
