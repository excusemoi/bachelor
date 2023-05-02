package model

import "time"

type TransformationRule struct {
	tableName   struct{}  `json:"-" pg:"transformation_rule,discard_unknown_columns"`
	Id          int       `json:"id" pg:"id,pk"`
	Field       string    `json:"field" pg:"field"`
	TargetField string    `json:"targetField" pg:"target_field"`
	UpdatedAt   time.Time `json:"updatedAt" pg:"updated_at, default:now()"`
}
