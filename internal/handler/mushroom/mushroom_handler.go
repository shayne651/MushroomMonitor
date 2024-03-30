package mushroom_handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	mushroom_service "github.com/shayne651/MushroomMonitor/internal/services/mushroom"
)

type IMushroomHandler interface {
	InitializeRestAPI(mux *http.ServeMux)
}

type MushroomHandler struct {
	MushroomService mushroom_service.IMushroomService
}

func (mh *MushroomHandler) InitializeRestAPI(mux *http.ServeMux) {
	log.Println("Initializing REST API for mushroom")
	mux.HandleFunc("POST /mushroom", mh.handleSaveMushroom)
	mux.HandleFunc("GET /mushroom", mh.handleGetMushrooms)
}

func (mh *MushroomHandler) handleSaveMushroom(w http.ResponseWriter, r *http.Request) {
	m := mushroom_service.Mushroom{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body from POST /mushroom", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println("Error unmarshaling body from POST /mushroom", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = mh.MushroomService.SaveMushroom(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (mh *MushroomHandler) handleGetMushrooms(w http.ResponseWriter, r *http.Request) {
	mushrooms, err := mh.MushroomService.GetMushrooms()
	if err != nil {
		log.Println("Error getting mushrooms", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	JSON, err := json.MarshalIndent(mushrooms, "", "\t")
	if err != nil {
		log.Println("Error marshling mushroom json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(JSON)
}
