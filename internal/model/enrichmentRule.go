package model

import "time"

type EnrichmentRule struct {
	tableName struct{}  `json:"-" pg:"enrichment_rule,discard_unknown_columns"`
	Id        int       `json:"id" pg:"id,pk"`
	Field     string    `json:"field" pg:"field"`
	Value     string    `json:"value" pg:"value"`
	UpdatedAt time.Time `json:"updatedAt" pg:"updated_at, default:now()"`
}
