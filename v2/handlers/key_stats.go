package handlers

import (
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/store"
	log "github.com/daiLlew/funkylog"
	"github.com/gorilla/mux"
	"net/http"
)

// GetProfileStatsHandlerFunc HTTP handler returns the current key stats for the specified area profile.
func GetProfileStatsHandlerFunc(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("handling %s request", "GET /profiles/{area_code}/stats")

		areaCode := mux.Vars(r)["area_code"]
		if areaCode == "" {
			http.Error(w, "area code required but none provided", http.StatusBadRequest)
			return
		}

		profile, err := db.GetProfileByAreaCode(areaCode)
		if err != nil {
			if err == store.ErrNotFound {
				http.Error(w, "profile not found", http.StatusNotFound)
				return
			}

			log.Err("error querying for profile: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		stats, err := db.GetKeyStatsForProfile(profile.ID)
		if err != nil {
			log.Err("error querying for profile stats: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		if err := writeEntity(w, stats, http.StatusOK); err != nil {
			log.Err("error writting get stats entity to reponse: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
