package mushroom_service

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type IMushroomService interface {
	GetMushrooms() ([]Mushroom, error)
	SaveMushroom(m Mushroom) error
}

type MushroomService struct {
	DB *sql.DB
}

type Mushroom struct {
	MushroomID uuid.UUID
	Name       string
}

func (ms *MushroomService) GetMushrooms() ([]Mushroom, error) {
	rows, err := ms.DB.Query("SELECT mushroom_uuid, name FROM mushroom")
	if err != nil {
		log.Println("error getting mushrooms", err)
		return nil, err

	}
	var mushrooms []Mushroom
	for rows.Next() {
		var mushroom Mushroom
		rows.Scan(&mushroom.MushroomID, &mushroom.Name)
		mushrooms = append(mushrooms, mushroom)
	}
	return mushrooms, nil
}

func (ms *MushroomService) SaveMushroom(m Mushroom) error {
	_, err := ms.DB.Exec("INSERT INTO mushroom(mushroom_uuid, name) VALUES (?, ?)", m.MushroomID, m.Name)
	if err != nil {
		log.Println("error inserting mushroom", err)
		return err
	}
	return nil
}
