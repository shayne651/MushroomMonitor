package sensor_service

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Fea struct {
	FeaID      uuid.UUID `json:"feaID"`
	GrowName   string    `json:"growName"`
	Runtime    float32   `json:"runtime"`
	RecordDate int       `json:"recordDate"`
}

type IFeaService interface {
	SaveFea(fea Fea) error
	GetLastFea(growName string) (Fea, error)
	GetRecentFea(days int, growName string) ([]Fea, error)
}

type FeaService struct {
	DB *sql.DB
}

func (fs *FeaService) GetLastFea(growName string) (Fea, error) {
	row := fs.DB.QueryRow("SELECT * from fea WHERE grow_name = ? ORDER BY record_date DESC LIMIT 1", growName)
	lastFEA := Fea{}
	err := row.Scan(&lastFEA.FeaID, &lastFEA.GrowName, &lastFEA.Runtime, &lastFEA.RecordDate)
	if err != nil {
		log.Println("Error getting last FEA for ", growName, ": ", err)
		return Fea{}, err
	}
	return lastFEA, nil
}

func (fs *FeaService) GetRecentFea(days int, growName string) ([]Fea, error) {
	var feas []Fea
	t := time.Now().AddDate(0, 0, -days)
	rows, err := fs.DB.Query("SELECT * FROM fea WHERE grow_name = ? AND record_date >= ?", growName, t)
	if err != nil {
		log.Println("Error getting FEA for range of ", days, "for grow ", growName, " :", err)
		return nil, err
	}
	for rows.Next() {
		fea := Fea{}
		err := rows.Scan(&fea.FeaID, &fea.GrowName, &fea.Runtime, &fea.RecordDate)
		if err != nil {
			log.Println("Error getting FEA for range of ", days, "for grow ", growName, " :", err)
			return nil, err
		}
		feas = append(feas, fea)
	}
	return feas, nil
}

func (fs *FeaService) SaveFea(fea Fea) error {
	fea.FeaID = uuid.New()
	fea.RecordDate = int(time.Now().Unix())
	_, err := fs.DB.Exec("INSERT INTO fea(fea_uuid, grow_name, runtime, record_date) VALUES(?, ?, ?, ?)", fea.FeaID, fea.GrowName, fea.Runtime, fea.RecordDate)
	if err != nil {
		log.Println("Error saving fea")
	}
	return nil
}
