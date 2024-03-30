package stage_handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	stage_service "github.com/shayne651/MushroomMonitor/internal/services/stage"
)

type IStageHandler interface {
}

type StageHandler struct {
	StageService stage_service.IStageService
}

func (sh *StageHandler) Initialize(mux *http.ServeMux) {
	mux.HandleFunc("GET /stage/{id}", sh.handleGetStage)
	mux.HandleFunc("POST /stage", sh.handleSaveStage)
}

func (sh *StageHandler) handleGetStage(w http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		log.Println("error parsing uuid", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid mushroom UUID"))
	}

	stages, err := sh.StageService.GetStagesForMushroom(parsedUuid)
	if err != nil {
		log.Println("error getting mushrooms", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.MarshalIndent(stages, "", "\t")
	stagesJson, err := json.MarshalIndent(stages, "", "\t")
	if err != nil {
		log.Println("error mashaling stages", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(stagesJson)
}

func (sh *StageHandler) handleSaveStage(w http.ResponseWriter, request *http.Request) {
	var stage stage_service.Stage
	json.NewDecoder(request.Body).Decode(&stage)
	err := sh.StageService.SaveStage(stage)
	if err != nil {
		log.Println("Error saving stage", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
