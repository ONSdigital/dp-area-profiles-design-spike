package store

import (
	"time"
)

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
	ProfileID    int                  `json:"profile_id"`
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
