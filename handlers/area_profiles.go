package handlers

import (
	"fmt"
	"github.com/ONSdigital/dp-area-profiles-design-spike/store"
	log "github.com/daiLlew/funkylog"
	"github.com/gorilla/mux"
	"net/http"
)

// GetAreaProfilesHandlerFunc http.HandlerFun returns a list available area profiles.
func GetAreaProfilesHandlerFunc(profilesStore store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profiles, err := profilesStore.GetProfiles()
		if err != nil {
			log.Err("GetAreaPorfilesHandler error %+v", err)
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
		}

		if err := writeEntity(w, profiles, http.StatusOK); err != nil {
			log.Err("GetAreaProfileHandler error %+v", err)
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
	}
}

// GetAreaProfileHandlerFunc http.HandlerFunc returns the area profile for the specified area code.
func GetAreaProfileHandlerFunc(profileStore store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		areaCode := mux.Vars(r)["area_code"]
		if areaCode == "" {
			http.Error(w, "area code required", http.StatusBadRequest)
			return
		}

		log.Info("handling getProfile request for areaCode=%s", areaCode)

		profile, err := profileStore.GetProfileByAreaCode(areaCode)
		if err != nil {
			if err == store.ErrNotFound {
				http.Error(w, "area code not found", http.StatusNotFound)
				return
			}
			log.Err("GetAreaProfileHandler error %+v", err)
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		keyStats, err := profileStore.GetKeyStatsByProfileID(profile.ID)
		if err != nil {
			log.Err("GetAreaProfileHandler error %+v", err)
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		profile.KeyStats = keyStats

		if err := writeEntity(w, profile, http.StatusOK); err != nil {
			log.Err("GetAreaProfileHandler error %+v", err)
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
	}
}
