package grow_service

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type IGrowService interface {
	getAllGrows()
}

type GrowService struct {
	DB *sql.DB
}

type Grow struct {
	Name       string
	Mushroom   uuid.UUID
	Stage      uuid.UUID
	Automation uuid.UUID
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
