package server

import (
	"Lighthouse/internal/database/gorm_pg_adapter"
	"Lighthouse/internal/database/spec/interfaces"
	"Lighthouse/internal/server/handlers/api"
	"Lighthouse/internal/server/middleware"
	"fmt"
	"net/http"
	"os"
	"time"
)

func createApplicationDb() interfaces.DatabaseConnectorStrategy {

	adapter, err := gorm_pg_adapter.CreateGormPgAdapter(
		"localhost",
		"postgre",
		"password",
		9001,
		"lighthouse_dev",
	)
	if err != nil {
		panic(err)
	}

	return adapter
}

func Serve() {

	port := os.Getenv("LH_PORT")
	if port == "" {
		panic("[LH_PORT] Port for server not set!")
	}

	mainRouter := http.NewServeMux()
	mainRouter.
		Handle(
			"/api/",
			http.StripPrefix(
				"/api",
				middleware.Apply(
					api.Router(),
					middleware.ApplyTimeout(time.Second*5),
					middleware.ApplyAttachDb(createApplicationDb()),
				),
			),
		)

	fmt.Printf("Starting server on %v\n", port)
	_ = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), mainRouter)
}
