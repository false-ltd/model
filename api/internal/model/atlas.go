package model

// AtlasData is the compact point cloud for the interactive model atlas.
// Field names are abbreviated to keep the payload small (~4k points).
type AtlasData struct {
	Providers []AtlasProvider `json:"providers"`
	Points    []AtlasPoint    `json:"points"`
}

type AtlasProvider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AtlasPoint struct {
	ID         uint     `json:"id"`
	N          string   `json:"n"` // display name
	M          string   `json:"m"` // model_id
	P          string   `json:"p"` // provider_id (index resolved client-side)
	CI         *float64 `json:"ci"` // input cost per 1M tokens
	CO         *float64 `json:"co"` // output cost
	CTX        *int     `json:"ctx"` // context window
	K          int      `json:"k"` // capability count (0-6)
	OW         bool     `json:"ow"` // open weights
}

// RecentModel is a homepage "recently indexed" entry.
type RecentModel struct {
	ID          uint   `gorm:"column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	ReleaseDate string `gorm:"column:release_date" json:"release_date"`
}
