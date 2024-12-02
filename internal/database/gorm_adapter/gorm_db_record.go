package gorm_adapter

import (
	models "Lighthouse/internal/database/models"
	"net/url"
	"time"
)

/* //////////////////////////////
Database struct and related records
////////////////////////////// */

// GormRecord Type used for DB representation
type GormRecord struct {
	Target    string
	Id        string    `gorm:"primaryKey"`
	createdAt time.Time `gorm:"autoCreateTime"`
	updatedAt time.Time `gorm:"autoUpdateTime"`
}

// ToRecord Method to convert GormRecord to Record
func (r *GormRecord) ToRecord() (models.Record, error) {
	targetUrl, urlErr := url.Parse(r.Target)
	if urlErr != nil {
		return models.Record{}, urlErr
	}
	return models.CreateRecord(*targetUrl, r.Id)
}

/* //////////////////////////////
Convert Record to GormRecord
////////////////////////////// */

func convertRecordToDbRecord(record models.Record) (GormRecord, error) {
	return GormRecord{
		Target: record.Target.String(),
		Id:     record.Id,
	}, nil
}
