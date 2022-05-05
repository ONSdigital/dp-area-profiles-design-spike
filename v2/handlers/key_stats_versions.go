package handlers

import (
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/store"
	log "github.com/daiLlew/funkylog"
	"github.com/gorilla/mux"
	"net/http"
)

// GetStatsVersionsHandlerFunc HTTP handler func returns a list of available key status versions for the specified area code.
func GetStatsVersionsHandlerFunc(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("handling %s request", "GET /profiles/{area_code}/stats/versions")

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

		versionsList, err := db.GetKeyStatsVersionsForProfile(profile)
		if err != nil {
			log.Err("error querying for profile stats versions: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		versions := store.KeyStatisticVersions{
			AreaProfile: *profile,
			Versions:    versionsList,
		}

		if err := writeEntity(w, versions, http.StatusOK); err != nil {
			log.Err("error writting get versions entity to response: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

// GetStatsVersionHandlerFunc HTTP handler func that returns key stats belonging to the specified version of an area profile.
func GetStatsVersionHandlerFunc(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("handling %s request", "GET /profiles/{area_code}/stats/versions/{version}")

		areaCode := mux.Vars(r)["area_code"]
		if areaCode == "" {
			http.Error(w, "area code required but none provided", http.StatusBadRequest)
			return
		}

		version := mux.Vars(r)["version"]
		if version == "" {
			http.Error(w, "version required but none provided", http.StatusBadRequest)
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

		stats, err := db.GetKeyStatsVersion(profile, version)
		if err != nil {
			log.Err("error querying for stats version: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		if err := writeEntity(w, stats, http.StatusOK); err != nil {
			log.Err("error writing get version entity to response: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
