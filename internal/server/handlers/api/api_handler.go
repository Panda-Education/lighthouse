package api

import (
	"net/http"
)

func Router() *http.ServeMux {

	router := http.NewServeMux()

	router.Handle("POST /insert", http.HandlerFunc(insertRecord))
	router.Handle("POST /update", http.HandlerFunc(updateRecord))
	router.Handle("POST /delete", http.HandlerFunc(deleteRecord))

	return router

}
