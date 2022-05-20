package handlers

import (
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/store"
	log "github.com/daiLlew/funkylog"
	"github.com/gorilla/mux"
	"net/http"
)

// GetAreaProfilesHandlerFunc http handler returning a list of all area profiles.
func GetAreaProfilesHandlerFunc(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("handling %s request", "GET /profiles")

		profiles, err := db.GetAreaProfiles()
		if err != nil {
			http.Error(w, "error getting area profiles list", http.StatusInternalServerError)
			return
		}

		if err := writeEntity(w, profiles, http.StatusOK); err != nil {
			log.Err("error writing getProfiles entity to response: %s", err.Error())
			http.Error(w, "error writting get area profiles enity to response", http.StatusInternalServerError)
		}
	}
}

// GetAreaProfile http handler returning the area profile associated with the provided area code.
func GetAreaProfileHandlerFunc(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("handling %s request", "GET /profiles/{area_code}")

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

		if err := writeEntity(w, profile, http.StatusOK); err != nil {
			log.Err("error writing get profile entity to response: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
