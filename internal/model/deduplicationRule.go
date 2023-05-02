package model

import "time"

type DeduplicationRule struct {
	tableName struct{}            `json:"-" pg:"deduplication_rule,discard_unknown_columns"`
	Id        int                 `json:"id" pg:"id,pk"`
	Field     string              `json:"field" pg:"field"`
	Value     map[string]struct{} `json:"value" pg:"-"`
	Start     time.Time           `json:"startTime" pg:"start_time"`
	End       time.Time           `json:"endTime" pg:"end_time"`
	UpdatedAt time.Time           `json:"updatedAt" pg:"updated_at, default:now()"`
}
