package server

import (
	"Lighthouse/internal/database/database_cache"
	"Lighthouse/internal/database/gorm_pg_adapter"
	"Lighthouse/internal/database/spec/interfaces"
	"Lighthouse/internal/server/handlers/api"
	"Lighthouse/internal/server/handlers/redirect"
	"Lighthouse/internal/server/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func createApplicationDb() interfaces.DatabaseConnectorStrategy {

	adapter, err := gorm_pg_adapter.CreateGormPgAdapter(
		"lh-pg",
		"postgres",
		"password",
		5432,
		"lighthouse_dev",
	)

	if err != nil {
		panic(err)
	}

	db := database_cache.CreateLruDb(adapter, 256)

	if err := db.Migrate(context.Background()); err != nil {
		panic(err)
	}

	return db
}

func Serve() {

	port := os.Getenv("LH_PORT")
	if port == "" {
		panic("[LH_PORT] Port for server not set!")
	}

	db := createApplicationDb()

	mainRouter := http.NewServeMux()
	mainRouter.
		Handle(
			"/api/",
			http.StripPrefix(
				"/api",
				middleware.Apply(
					api.Router(),
					middleware.ApplyTimeout(time.Second*5),
					middleware.ApplyAttachDb(db),
				),
			),
		)
	mainRouter.
		Handle(
			"/",
			middleware.Apply(
				redirect.Router(),
				middleware.ApplyTimeout(time.Second*2),
				middleware.ApplyAttachDb(db),
			),
		)

	fmt.Printf("Starting server on %v\n", port)
	_ = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), mainRouter)
}
