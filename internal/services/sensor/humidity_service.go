package sensor_service

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Humidity struct {
	HumidityID uuid.UUID `json:"humidityID"`
	GrowName   string    `json:"growName"`
	Value      float32   `json:"value"`
	RecordDate int       `json:"recordDate"`
}

type IHumidityService interface {
	SaveHumidity(humidity Humidity) error
	GetLastHumidity(growName string) (Humidity, error)
	GetRecentHumidity(days int, growName string) ([]Humidity, error)
}

type HumidityService struct {
	DB *sql.DB
}

func (hs *HumidityService) GetLastHumidity(growName string) (Humidity, error) {
	row := hs.DB.QueryRow("SELECT * FROM humidity WHERE grow_name = ? ORDER BY record_date DESC LIMIT(1)", growName)
	humidity := Humidity{}
	err := row.Scan(&humidity.HumidityID, &humidity.GrowName, &humidity.Value, &humidity.RecordDate)
	if err != nil {
		log.Println("Error getting last humidity for ", growName, err)
		return humidity, err
	}
	return humidity, nil
}

func (hs *HumidityService) GetRecentHumidity(days int, growName string) ([]Humidity, error) {
	var humidities []Humidity
	t := time.Now().AddDate(0, 0, -days)
	rows, err := hs.DB.Query("SELECT * FROM humidity WHERE grow_name = ? AND record_date >= ?", growName, t)
	if err != nil {
		log.Println("Error getting Humidity for range of ", days, "for grow ", growName, " :", err)
		return nil, err
	}
	for rows.Next() {
		humidity := Humidity{}
		err := rows.Scan(&humidity.HumidityID, &humidity.GrowName, &humidity.Value, &humidity.RecordDate)
		if err != nil {
			log.Println("Error getting Humidity for range of ", days, "for grow ", growName, " :", err)
			return nil, err
		}
		humidities = append(humidities, humidity)
	}
	return humidities, nil
}

func (hs *HumidityService) SaveHumidity(humidity Humidity) error {
	humidity.HumidityID = uuid.New()
	humidity.RecordDate = int(time.Now().Unix())
	_, err := hs.DB.Exec("INSERT INTO humidity(humidity_uuid, grow_name, value, record_date) VALUES(?, ?, ?, ?)", humidity.HumidityID, humidity.GrowName, humidity.Value, humidity.RecordDate)
	if err != nil {
		log.Println("Error saving humidity")
	}
	return nil
}
