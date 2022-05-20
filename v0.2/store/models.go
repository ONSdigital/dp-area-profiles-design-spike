package store

import (
	"time"
)

// KetStatType provides a unique identity of each type of key stat value.
type KeyStatType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AreaProfile is a domain representation of a geographical area profile.
type AreaProfile struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	AreaCode string `json:"area_code"`
	Href     string `json:"href"`
}

type KeyStatistics []KeyStatistic

// KeyStatistic is a domain model representing a key statistical figure for an area profile.
type KeyStatistic struct {
	VersionID    int                  `json:"version_id,,omitempty"`
	StatID       int                  `json:"id"`
	StatType     int                  `json:"stat_type"`
	ProfileID    int                  `json:"-"`
	AreaCode     string               `json:"area_code"`
	Name         string               `json:"name"`
	Value        string               `json:"value"`
	Unit         string               `json:"unit"`
	DateCreated  time.Time            `json:"date_created"`
	LastModified time.Time            `json:"last_modified,omitempty"`
	Metadata     KeyStatisticMetadata `json:"metadata,omitempty"`
}

// KeyStatisticMetadata is a domain model representing metadata associated with a KeyStatistic
type KeyStatisticMetadata struct {
	DatasetID   string `json:"dataset_id"`
	DatasetName string `json:"dataset_name"`
	Href        string `json:"href"`
}

type KeyStatisticVersions struct {
	AreaProfile
	Versions []time.Time `json:"versions"`
}

// KeyStatsRecipe is a type encapslating the import job for a new dataset version notification.
// It specifies what Cantabular query to run, which geographies it affected too and which key stat type the query results represent.
type KeyStatsRecipe struct {
	// Unique ID for the recipe
	ID int
	// The Dataset ID the recipe allies to.
	DatasetID string
	// The Dataset Edition the recipe applies to.
	DatasetEdition string
	// A Cantabular query template to execute for this recipe.
	CantabularQuery string
	// The Key Stat type the query results represent to i.e. Resident Population
	StatType int
	// The Geography types this recipe applies to.
	Geographies []GeographyType
}

type GeographyType struct {
	// ID is an internal ID to uniquely identify a geography
	ID int
	// Code is the hierarchy code assigned to this geography type.
	Code string
	// The display name of the geography.
	Name string
}
