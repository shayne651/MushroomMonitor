package sensor_service

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Humidity struct {
	HumidityID uuid.UUID `json:"humidityID"`
	GrowID     uuid.UUID `json:"growID"`
	Value      float32   `json:"value"`
	RecordDate int       `json:"recordDate"`
}

type IHumidityService interface {
	SaveHumidity(humidity Humidity) error
	GetLastHumidity(growID uuid.UUID) (Humidity, error)
	GetRecentHumidity(days int, growID uuid.UUID) ([]Humidity, error)
}

type HumidityService struct {
	DB *sql.DB
}

func (hs *HumidityService) GetLastHumidity(growID uuid.UUID) (Humidity, error) {
	row := hs.DB.QueryRow("SELECT * FROM humidity WHERE grow_uuid = ? ORDER BY record_date DESC LIMIT(1)", growID)
	humidity := Humidity{}
	err := row.Scan(&humidity.HumidityID, &humidity.GrowID, &humidity.Value, &humidity.RecordDate)
	if err != nil {
		log.Println("Error getting last humidity for ", growID, err)
		return humidity, err
	}
	return humidity, nil
}

func (hs *HumidityService) GetRecentHumidity(days int, growID uuid.UUID) ([]Humidity, error) {
	var humidities []Humidity
	t := time.Now().AddDate(0, 0, -days)
	rows, err := hs.DB.Query("SELECT * FROM humidity WHERE grow_uuid = ? AND record_date >= ?", growID, t)
	if err != nil {
		log.Println("Error getting Humidity for range of ", days, "for grow ", growID, " :", err)
		return nil, err
	}
	for rows.Next() {
		humidity := Humidity{}
		err := rows.Scan(&humidity.HumidityID, &humidity.GrowID, &humidity.Value, &humidity.RecordDate)
		if err != nil {
			log.Println("Error getting Humidity for range of ", days, "for grow ", growID, " :", err)
			return nil, err
		}
		humidities = append(humidities, humidity)
	}
	return humidities, nil
}

func (hs *HumidityService) SaveHumidity(humidity Humidity) error {
	humidity.HumidityID = uuid.New()
	humidity.RecordDate = int(time.Now().Unix())
	_, err := hs.DB.Exec("INSERT INTO humidity(humidity_uuid, grow_uuid, value, record_date) VALUES(?, ?, ?, ?)", humidity.HumidityID, humidity.GrowID, humidity.Value, humidity.RecordDate)
	if err != nil {
		log.Println("Error saving humidity")
	}
	return nil
}
