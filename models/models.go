package models

import (
	"fmt"
	"time"
)

// ImportRow is a Go representation of an area profiles key statistic in an import cvs row.
type ImportRow struct {
	AreaCode string
	Title    string
	Name     string
	Value    string
	Unit     string
}

// Area is a simplified domain representation of a Geographical area. Other fields omitted for the sake of this POC.
type Area struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// AreaProfileLink is a model encapsulating details to link to an area profile
type AreaProfileLink struct {
	ProfileID int    `json:"profile_id"`
	Name      string `json:"name"`
	AreaCode  string `json:"area_code"`
	Href      string `json:"href"`
}

// AreaProfile is a domain representation of a geographical area profile.
type AreaProfile struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	AreaCode string         `json:"area_code"`
	KeyStats []KeyStatistic `json:"key_stats"`
}

// KeyStatistic is a domain model representing a key statistical figure for an area profile.
type KeyStatistic struct {
	VersionID    int       `json:"version_id,,omitempty"`
	StatID       int       `json:"id"`
	ProfileID    int       `json:"profile_id"`
	Name         string    `json:"name"`
	Value        string    `json:"value"`
	Unit         string    `json:"unit"`
	DateCreated  time.Time `json:"date_created"`
	LastModified time.Time `json:"last_modified,omitempty"`
}

// KeyStatsVersion is model containing verison metda data about a an area profile key statistic.
type KeyStatsVersion struct {
	StatID       int       `json:"id"`
	ProfileID    int       `json:"profile_id"`
	VersionID    int       `json:"version_id"`
	DateCreated  time.Time `json:"date_created"`
	LastModified time.Time `json:"last_modified"`
	Href         string    `json:"href"`
}

// ToString is a function that produces a string representation of KeyStatistic.
func (k KeyStatistic) ToString() string {
	format := "[StatID: %d, ProfileID: %d, Name: %s, Value: %s, Unit: %s, Date Created: %+v]"
	return fmt.Sprintf(format, k.StatID, k.ProfileID, k.Name, k.Value, k.Unit, k.DateCreated)
}
