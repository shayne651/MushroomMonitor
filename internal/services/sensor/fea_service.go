package sensor_service

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Fea struct {
	FeaID      uuid.UUID `json:"feaID"`
	GrowID     uuid.UUID `json:"growID"`
	Runtime    float32   `json:"runtime"`
	RecordDate int       `json:"recordDate"`
}

type IFeaService interface {
	SaveFea(fea Fea) error
	GetLastFea(growID uuid.UUID) (Fea, error)
	GetRecentFea(days int, growID uuid.UUID) ([]Fea, error)
}

type FeaService struct {
	DB *sql.DB
}

func (fs *FeaService) GetLastFea(growID uuid.UUID) (Fea, error) {
	row := fs.DB.QueryRow("SELECT * from fea WHERE grow_uuid = ? ORDER BY record_date DESC LIMIT 1", growID)
	lastFEA := Fea{}
	err := row.Scan(&lastFEA.FeaID, &lastFEA.GrowID, &lastFEA.Runtime, &lastFEA.RecordDate)
	if err != nil {
		log.Println("Error getting last FEA for ", growID, ": ", err)
		return Fea{}, err
	}
	return lastFEA, nil
}

func (fs *FeaService) GetRecentFea(days int, growID uuid.UUID) ([]Fea, error) {
	var feas []Fea
	t := time.Now().AddDate(0, 0, -days)
	rows, err := fs.DB.Query("SELECT * FROM fea WHERE grow_uuid = ? AND record_date >= ?", growID, t)
	if err != nil {
		log.Println("Error getting FEA for range of ", days, "for grow ", growID, " :", err)
		return nil, err
	}
	for rows.Next() {
		fea := Fea{}
		err := rows.Scan(&fea.FeaID, &fea.GrowID, &fea.Runtime, &fea.RecordDate)
		if err != nil {
			log.Println("Error getting FEA for range of ", days, "for grow ", growID, " :", err)
			return nil, err
		}
		feas = append(feas, fea)
	}
	return feas, nil
}

func (fs *FeaService) SaveFea(fea Fea) error {
	fea.FeaID = uuid.New()
	fea.RecordDate = int(time.Now().Unix())
	_, err := fs.DB.Exec("INSERT INTO fea(fea_uuid, grow_uuid, runtime, record_date) VALUES(?, ?, ?, ?)", fea.FeaID, fea.GrowID, fea.Runtime, fea.RecordDate)
	if err != nil {
		log.Println("Error saving fea")
	}
	return nil
}
