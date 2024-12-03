package gorm_pg_adapter

import (
	models "Lighthouse/internal/database/models"
	"net/url"
	"time"
)

/* //////////////////////////////
Database struct and related records
////////////////////////////// */

// GormPgRecord Type used for DB representation
type GormPgRecord struct {
	Target    string
	Id        string    `gorm:"primaryKey"`
	createdAt time.Time `gorm:"autoCreateTime"`
	updatedAt time.Time `gorm:"autoUpdateTime"`
}

// ToRecord Method to convert GormPgRecord to Record
func (r *GormPgRecord) ToRecord() (models.Record, error) {
	targetUrl, urlErr := url.Parse(r.Target)
	if urlErr != nil {
		return models.Record{}, urlErr
	}
	return models.CreateRecord(targetUrl, r.Id)
}

/* //////////////////////////////
Convert Record to GormPgRecord
////////////////////////////// */

func convertRecordToDbRecord(record models.Record) (GormPgRecord, error) {
	return GormPgRecord{
		Target: record.Target.String(),
		Id:     record.Id,
	}, nil
}
