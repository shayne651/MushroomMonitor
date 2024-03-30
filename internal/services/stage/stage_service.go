package stage_service

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type IStageService interface {
  SaveStage(stage Stage) error
  GetStagesForMushroom(mushroomID uuid.UUID) ([]Stage, error)
  GetStageByPk(stageID, MushroomID uuid.UUID) (Stage, error)
}

type StageService struct {
	DB *sql.DB
}

type Stage struct {
	StageID     uuid.UUID
	MushroomID  uuid.UUID
	Name        string
	MinTemp     float32
	MaxTemp     float32
	MinHumidity float32
	MaxHumidity float32
	Fea         float32
}

func (ss *StageService) GetStagesForMushroom(mushroomID uuid.UUID) ([]Stage, error) {
	rows, err := ss.DB.Query("SELECT stage_uuid, mushroom_uuid, name, min_temp, max_temp, min_humidity, max_humidity, fea FROM stage WHERE mushroom_uuid = ?", mushroomID)
	if err != nil {
		log.Println("Error getting stages for mushroom", err)
		return nil, err
	}

	var stages []Stage

	for rows.Next() {
		var s Stage
		rows.Scan(&s.StageID, &s.MushroomID, &s.Name, &s.MinTemp, &s.MaxTemp, &s.MinHumidity, &s.MaxHumidity, &s.Fea)
		stages = append(stages, s)
	}
	return stages, nil
}

func (ss *StageService) GetStageByPk(stageID, mushroomID uuid.UUID) (Stage, error) {
	row := ss.DB.QueryRow("SELECT stage_uuid, mushroom_uuid, name, min_temp, max_temp, min_humidity, max_humidity, fea FROM stage WHERE mushroom_uuid = ? AND stage_id = ?", mushroomID, stageID)

	var s Stage

	err := row.Scan(&s.StageID, &s.MushroomID, &s.Name, &s.MinTemp, &s.MaxTemp, &s.MinHumidity, &s.MaxHumidity, &s.Fea)
	if err != nil {
		log.Println("Error getting stage by PK", err)
		return Stage{}, err
	}
	return s, nil

}

func (ss *StageService) SaveStage(stage Stage) error {
	if stage.StageID != uuid.Nil {
    var err error
    stage.StageID, err = uuid.NewRandom()
		if err != nil {
			log.Println("Failed to generate new uuid for stage")
			return err
		}
	}
	_, err := ss.DB.Exec("INSERT INTO stage (stage_uuid, mushroom_uuid, name, min_temp, max_temp, min_humidity, max_humidity, fea) VALUES (?,?,?,?,?,?,?,?)",
		stage.StageID, stage.MushroomID, stage.Name, stage.MinTemp, stage.MaxTemp, stage.MinHumidity, stage.MaxHumidity, stage.Fea)
  if err != nil {
    log.Println("Error saving new stage ", err)
    return err
  }
  return nil
}
