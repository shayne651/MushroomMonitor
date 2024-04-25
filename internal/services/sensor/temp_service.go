package sensor_service

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Temp struct {
	TempID     uuid.UUID `json:"tempID"`
	GrowName   string    `json:"growName"`
	Value      float32   `json:"value"`
	RecordDate int       `json:"recordDate"`
}

type ITempService interface {
	SaveTemp(temp Temp) error
	GetLastTemp(name string) (Temp, error)
	GetRecentTemp(days int, name string) ([]Temp, error)
}

type TempService struct {
	DB *sql.DB
}

func (ts *TempService) GetLastTemp(name string) (Temp, error) {
	row := ts.DB.QueryRow("SELECT * from temp WHERE grow_name = ? ORDER BY record_date DESC LIMIT 1", name)
	lastTemp := Temp{}
	err := row.Scan(&lastTemp.TempID, &lastTemp.GrowName, &lastTemp.Value, &lastTemp.RecordDate)
	if err != nil {
		log.Println("Error getting last temp for ", name, ": ", err)
		return Temp{}, err
	}
	return lastTemp, nil
}

func (ts *TempService) GetRecentTemp(days int, name string) ([]Temp, error) {
	var temps []Temp
	t := time.Now().AddDate(0, 0, -days)
	rows, err := ts.DB.Query("SELECT * FROM temp WHERE grow_name = ? AND record_date >= ?", name, t)
	if err != nil {
		log.Println("Error getting temp for range of ", days, "for grow ", name, " :", err)
		return nil, err
	}
	for rows.Next() {
		temp := Temp{}
		err := rows.Scan(&temp.TempID, &temp.GrowName, &temp.Value, &temp.RecordDate)
		if err != nil {
			log.Println("Error getting temp for range of ", days, "for grow ", name, " :", err)
			return nil, err
		}
		temps = append(temps, temp)
	}
	return temps, nil
}

func (ts *TempService) SaveTemp(temp Temp) error {
	temp.TempID = uuid.New()
	temp.RecordDate = int(time.Now().Unix())
	_, err := ts.DB.Exec("INSERT INTO temp(temp_uuid, grow_name, value, record_date) VALUES(?, ?, ?, ?)", temp.TempID, temp.GrowName, temp.Value, temp.RecordDate)
	if err != nil {
		log.Println("Error saving temp", err)
	}
	return nil
}
