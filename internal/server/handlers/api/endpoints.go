package api

import (
	"Lighthouse/internal/server/io"
	"Lighthouse/internal/server/middleware"
	"encoding/json"
	"net/http"
)

func insertRecord(w http.ResponseWriter, r *http.Request) {
	db, ok := middleware.GetDbFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to load DB instance from middleware", http.StatusInternalServerError)
		return
	}

	var recordIO io.RecordIO
	if err := json.NewDecoder(r.Body).Decode(&recordIO); err != nil {
		http.Error(w, "Unable to parse RecordIO request", http.StatusBadRequest)
		return
	}

	record, _ := recordIO.ToRecord()

	if err := db.InsertRecord(r.Context(), record); err != nil {
		http.Error(w, "Unable to insert Record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(io.StatusIO{Message: "Success"}); err != nil {
		http.Error(w, "Unable to create response message", http.StatusInternalServerError)
	}
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	db, ok := middleware.GetDbFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to load DB instance from middleware", http.StatusInternalServerError)
		return
	}

	var recordIO io.RecordIO
	if err := json.NewDecoder(r.Body).Decode(&recordIO); err != nil {
		http.Error(w, "Unable to parse RecordIO request", http.StatusBadRequest)
		return
	}

	record, _ := recordIO.ToRecord()

	if err := db.UpdateRecord(r.Context(), record); err != nil {
		http.Error(w, "Unable to update Record, does it exist?", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(io.StatusIO{Message: "Success"}); err != nil {
		http.Error(w, "Unable to create response message", http.StatusInternalServerError)
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	db, ok := middleware.GetDbFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to load DB instance from middleware", http.StatusInternalServerError)
		return
	}

	var recordIO io.RecordIO
	if err := json.NewDecoder(r.Body).Decode(&recordIO); err != nil {
		http.Error(w, "Unable to parse RecordIO request", http.StatusBadRequest)
		return
	}

	record, _ := recordIO.ToRecord()

	if err := db.DeleteRecord(r.Context(), record.Id); err != nil {
		http.Error(w, "Unable to delete Record, does it exist?", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(io.StatusIO{Message: "Success"}); err != nil {
		http.Error(w, "Unable to create response message", http.StatusInternalServerError)
	}
}
