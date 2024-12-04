package redirect

import "net/http"

func Router() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("GET /{target}", http.HandlerFunc(index))
	return router
}
