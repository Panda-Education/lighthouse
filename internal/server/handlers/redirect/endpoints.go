package redirect

import (
	"Lighthouse/internal/server/middleware"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {

	target := r.PathValue("target")
	if target == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	db, ok := middleware.GetDbFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to load DB instance from middleware", http.StatusInternalServerError)
		return
	}

	record, err := db.FindRecord(r.Context(), target)
	if err != nil {
		http.Error(w, "Unable to retrieve target from database", http.StatusInternalServerError)
		return
	}
	if record == nil {
		http.Error(w, "Target was not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, record.Target.String(), http.StatusFound)
}
