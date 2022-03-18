package handlers

import (
	"encoding/json"
	"github.com/ONSdigital/dp-area-profiles-design-spike/store"
	"github.com/gorilla/mux"
	"net/http"
)

//
func Initialise(s store.Store) (*mux.Router, error) {
	r := mux.NewRouter()
	r.Path("/profiles").Methods(http.MethodGet).HandlerFunc(GetAreaProfilesHandlerFunc(s))
	r.Path("/profiles/{area_code}").Methods(http.MethodGet).HandlerFunc(GetAreaProfileHandlerFunc(s))
	r.Path("/profiles/{area_code}/{file}").Methods(http.MethodPut).HandlerFunc(AddNewKeyStatsVersion(s))
	r.Path("/profiles/{area_code}/versions").Methods(http.MethodGet).HandlerFunc(GetKeyStatVersionsHandlerFunc(s))
	r.Path("/profiles/{area_code}/versions/{version_id}").Methods(http.MethodGet).HandlerFunc(GetKeyStatVersionHandlerFunc(s))

	return r, nil
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
