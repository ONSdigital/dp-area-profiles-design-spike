package handlers

import (
	"fmt"
	"github.com/ONSdigital/dp-area-profiles-design-spike/store"
	"github.com/ONSdigital/dp-area-profiles-design-spike/testdata"
	log "github.com/daiLlew/funkylog"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// AddNewKeyStatsVersion http.HandlerFunc imports a new version of the area profile key statistics test data.
// The "current" figures are versioned before the new figures are inserted.
func AddNewKeyStatsVersion(profileStore store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		areaCode := mux.Vars(r)["area_code"]
		importFile := mux.Vars(r)["file"]

		if areaCode == "" {
			http.Error(w, "area code required", http.StatusBadRequest)
			return
		}

		profile, err := profileStore.GetProfileByAreaCode(areaCode)
		if err != nil {
			if err == store.ErrNotFound {
				http.Error(w, "area code not found", http.StatusNotFound)
				return
			}
			log.Err(err.Error())
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		newStats, err := testdata.ReadCVS(fmt.Sprintf("testdata/%s.csv", importFile))
		if err != nil {
			log.Err(err.Error())
			http.Error(w, "error reading text data csv", http.StatusInternalServerError)
			return
		}

		if err := profileStore.UpdateProfileKeyStats(profile.ID, newStats); err != nil {
			log.Err(err.Error())
			http.Error(w, "error updating profile key stats", http.StatusInternalServerError)
			return
		}

		entity := map[string]string{
			"message":  "new profile key stats version created successfully",
			"href":     fmt.Sprintf("http://localhost:8080/profiles/%s", areaCode),
			"versions": fmt.Sprintf("http://localhost:8080/profiles/%s/versions", areaCode),
		}

		if err := writeEntity(w, entity, http.StatusOK); err != nil {
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
	}
}

// GetKeyStatVersionsHandlerFunc http.HandlerFunc returning a list of available key statistics versions for an area profile.
func GetKeyStatVersionsHandlerFunc(db store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		areaCode := mux.Vars(r)["area_code"]
		if areaCode == "" {
			http.Error(w, "area code required", http.StatusBadRequest)
			return
		}

		profile, err := db.GetProfileByAreaCode(areaCode)
		if err != nil {
			if err == store.ErrNotFound {
				http.Error(w, "area code not found", http.StatusNotFound)
				return
			}
			log.Err(err.Error())
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		versions, err := db.GetKeyStatsVersions(areaCode, profile.ID)
		if err != nil {
			log.Err(err.Error())
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if err := writeEntity(w, versions, http.StatusOK); err != nil {
			log.Err(err.Error())
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
	}
}

// GetKeyStatVersionHandlerFunc http.HandlerFunc returns key stats associated with a the specific versionID & profileID combination.
func GetKeyStatVersionHandlerFunc(db store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		areaCode := mux.Vars(r)["area_code"]
		versionIDStr := mux.Vars(r)["version_id"]

		versionID, err := strconv.Atoi(versionIDStr)
		if err != nil {
			http.Error(w, "invalid version id", http.StatusBadRequest)
			return
		}

		profile, err := db.GetProfileByAreaCode(areaCode)
		if err != nil {
			if err == store.ErrNotFound {
				http.Error(w, "area code not found", http.StatusNotFound)
				return
			}
			log.Err(err.Error())
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		stats, err := db.GetKeyStatsVersion(profile.ID, versionID)
		if err != nil {
			log.Err(err.Error())
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if err := writeEntity(w, stats, http.StatusOK); err != nil {
			log.Err(err.Error())
			http.Error(w, fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
	}
}
