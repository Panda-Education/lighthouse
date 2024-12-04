package io

import "Lighthouse/internal/database/models"

type RecordIO struct {
	Target string `json:"target"`
	Id     string `json:"id"`
	Uid    string `json:"uid"`
}

func (r *RecordIO) ToRecord() (*models.Record, error) {
	return models.CreateRecordFromString(r.Target, r.Id, r.Uid)
}
