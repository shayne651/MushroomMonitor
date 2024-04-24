package grow_service

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type IGrowService interface {
	GetAllGrows()
	GetFullGrow(id uuid.UUID) (Grow, error)
}

type GrowService struct {
	DB *sql.DB
}

type Grow struct {
	Name       string    `json:"name"`
	Mushroom   uuid.UUID `json:"mushroom"`
	Stage      uuid.UUID `json:"stage"`
	Automation uuid.UUID `json:"automation"`
}

type FullGrow struct {
	GrowName     string  `json:"growName"`
	MushroomName string  `json:"mushroomName"`
	StageName    string  `json:"stageName"`
	MinTemp      float32 `json:"minTemp"`
	MaxTemp      float32 `json:"maxTemp"`
	MinHumidity  float32 `json:"minHumidity"`
	MaxHumidity  float32 `json:"maxHumidity"`
	Fea          float32 `json:"fea"`
}

func (gs *GrowService) GetFullGrow(name string) (FullGrow, error) {
	row := gs.DB.QueryRow(
		"SELECT g.name, m.name, s.name, s.min_temp, s.max_temp, s.min_humidity, s.max_humidity, s.fea "+
			"FROM grow g "+
			"JOIN "+
			"mushroom m ON g.mushroom_uuid = m.mushroom_uuid "+
			"JOIN "+
			"stage s ON g.stage_uuid = s.stage_uuid AND s.mushroom_uuid = m.mushroom_uuid "+
			"WHERE g.name = ?", name)
	var g FullGrow
	err := row.Scan(&g.GrowName, &g.MushroomName, &g.StageName, &g.MinTemp, &g.MaxTemp, &g.MinHumidity, &g.MaxHumidity, g.Fea)
	if err != nil {
		log.Println("Error getting full grow for ", name, ": ", err)
		return FullGrow{}, err
	}
	return g, nil
}

func (gs *GrowService) GetAllGrows() ([]Grow, error) {
	rows, err := gs.DB.Query("SELECT name, mushroom_uuid, stage_uuid, automation_uuid FROM grow")
	if err != nil {
		log.Println("Failed to get grows", err)
		return nil, err
	}
	var grows []Grow
	for rows.Next() {
		grow := Grow{}
		rows.Scan(&grow.Name, &grow.Mushroom, &grow.Stage, &grow.Automation)
		grows = append(grows, grow)
	}
	log.Println(grows)
	return grows, nil
}

func (gs *GrowService) SaveGrow(grow Grow) error {
	_, err := gs.DB.Exec("INSERT INTO grow (name, mushroom_uuid, stage_uuid, automation_uuid) VALUES (?, ?, ? ,?)", grow.Name, grow.Mushroom, grow.Stage, grow.Automation)
	if err != nil {
		log.Println("Error creating grow", err)
		return err
	}
	return nil
}

func (gs *GrowService) UpdateGrow(grow Grow) error {
	_, err := gs.DB.Exec("UPDATE grow SET mushroom_uuid = ?, stage_uuid = ? WHERE name = ? AND automation_uuid = ?", grow.Mushroom, grow.Stage, grow.Name, grow.Automation)
	if err != nil {
		log.Println("Error updating grow: ", grow.Name, " ", err)
	}
	return nil
}

func (gs *GrowService) GetGrow(name string) Grow {
	row := gs.DB.QueryRow("SELECT name, mushroom_uuid, stage_uuid, automation_uuid FROM grow WHERE name = ?", name)
	var g Grow
	row.Scan(&g.Name, &g.Mushroom, &g.Stage, &g.Automation)
	return g
}
