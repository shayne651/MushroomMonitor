package grow_handler

import (
	grow_service "github.com/shayne651/MushroomMonitor/internal/services/grow"

	"encoding/json"
	"log"
	"net/http"
)

type IGrowHandler interface {
}

type GrowHandler struct {
	GrowService grow_service.GrowService
}

func (gh *GrowHandler) Initialize(mux *http.ServeMux) {
	mux.HandleFunc("GET /grows", gh.handleGetGrows)
	mux.HandleFunc("GET /grow/{name}", gh.handleGetGrow)
	mux.HandleFunc("POST /grow", gh.handleSaveGrow)
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
