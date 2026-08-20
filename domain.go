package address

import "time"

// Match is a resolved street with a confidence score in [0, 1].
type Match struct {
	// Name is the canonical street name as stored in the database.
	Name string `json:"name"`
	// NormalizedName is the lowercased form used for lookups.
	NormalizedName string `json:"normalized_name"`
	// Confidence is the match score; 1.0 for exact matches.
	Confidence float64 `json:"confidence"`
}

type streetRow struct {
	OsmID          int64
	NameNormalized string
	NameOriginal   string
	Region         string
	LastUpdated    time.Time
}
