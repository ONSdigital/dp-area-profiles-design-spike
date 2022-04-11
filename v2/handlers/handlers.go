package handlers

import (
	"encoding/json"
	"github.com/ONSdigital/dp-area-profiles-design-spike/v2/store"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// DB represents the area profiles data store.
type DB interface {
	GetAreaProfiles() ([]store.AreaProfile, error)
	GetProfileByAreaCode(areaCode string) (*store.AreaProfile, error)
	GetKeyStatsForProfile(profileID int) (store.KeyStatistics, error)
	GetKeyStatsVersionsForProfile(profileID int) ([]time.Time, error)
	GetKeyStatsVersion(profileID int, date string) (store.KeyStatistics, error)
}

// Initalise registers the API handler functions.
func Initalise(db DB) *mux.Router {
	r := mux.NewRouter()

	r.Path("/profiles").Methods(http.MethodGet).HandlerFunc(GetAreaProfilesHandlerFunc(db))
	r.Path("/profiles/{area_code}").Methods(http.MethodGet).HandlerFunc(GetAreaProfileHandlerFunc(db))
	r.Path("/profiles/{area_code}/stats").Methods(http.MethodGet).HandlerFunc(GetProfileStatsHandlerFunc(db))
	r.Path("/profiles/{area_code}/stats/versions").Methods(http.MethodGet).HandlerFunc(GetStatsVersionsHandlerFunc(db))
	r.Path("/profiles/{area_code}/stats/versions/{version}").Methods(http.MethodGet).HandlerFunc(GetStatsVersionHandlerFunc(db))
	return r
}

func writeEntity(w http.ResponseWriter, entity interface{}, status int) error {
	body, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		return err
	}

	return nil
}
