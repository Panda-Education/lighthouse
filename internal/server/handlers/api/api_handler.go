package api

import (
	"fmt"
	"net/http"
)

func Router() *http.ServeMux {

	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello world!")
	}))

	return router

}
