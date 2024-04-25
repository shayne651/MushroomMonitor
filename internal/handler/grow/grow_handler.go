package grow_handler

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	grow_service "github.com/shayne651/MushroomMonitor/internal/services/grow"
	message_service "github.com/shayne651/MushroomMonitor/internal/services/message"

	"encoding/json"
	"log"
	"net/http"
)

type IGrowHandler interface {
}

type GrowHandler struct {
	GrowService grow_service.GrowService
	MQ          *message_service.MqttService
}

func (gh *GrowHandler) Initialize(mux *http.ServeMux) {
	mux.HandleFunc("GET /grows", gh.handleGetGrows)
	mux.HandleFunc("GET /grow/{name}", gh.handleGetGrow)
	mux.HandleFunc("POST /grow", gh.handleSaveGrow)
	mux.HandleFunc("GET /grow/full/{name}", gh.handleGetFullGrow)

	gh.MQ.SubscribeToTopic("mushroom_monitor-test.request-config", gh.handleRequestConfig)
}

func (gh *GrowHandler) handleGetGrows(w http.ResponseWriter, request *http.Request) {
	grows, err := gh.GrowService.GetAllGrows()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	growsJson, err := json.Marshal(grows)
	if err != nil {
		log.Println("Error marshaling grows", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(growsJson)
}

func (gh *GrowHandler) handleGetGrow(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")
	grow := gh.GrowService.GetGrow(name)

	growJson, err := json.Marshal(grow)
	if err != nil {
		log.Println("Error marshaling grow: ", name, " ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(growJson)
}

func (gh *GrowHandler) handleGetFullGrow(w http.ResponseWriter, request *http.Request) {
	name := request.PathValue("name")
	grow, _ := gh.GrowService.GetFullGrow(name)

	growJson, err := json.Marshal(grow)
	if err != nil {
		log.Println("Error marshaling full grow: ", name, " ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(growJson)
}

func (gh *GrowHandler) handleSaveGrow(w http.ResponseWriter, request *http.Request) {
	var g grow_service.Grow
	json.NewDecoder(request.Body).Decode(&g)
	log.Println(g)
	err := gh.GrowService.SaveGrow(g)
	if err != nil {
		log.Println("Error saving grow", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (gh *GrowHandler) handleRequestConfig(client mqtt.Client, msg mqtt.Message) {
	// TODO: Proper error handling, message should contain the name of the grow that the config is being requested for
	gs, _ := gh.GrowService.GetFullGrow("test")
	growJson, err := json.Marshal(gs)
	if err != nil {
		log.Println("Error marshaling full grow: ", "test", " ", err)
		return
	}
	gh.MQ.PublishMessage("mushroom_monitor-test.config", growJson)
}
