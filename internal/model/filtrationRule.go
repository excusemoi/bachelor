package model

import (
	"time"
)

type FiltrationRule struct {
	tableName struct{}  `json:"-" pg:"filtration_rule,discard_unknown_columns"`
	Id        int       `json:"id" pg:"id,pk"`
	Field     string    `json:"field" pg:"field"`
	Function  string    `json:"function" pg:"function"`
	Value     string    `json:"value" pg:"value"`
	UpdatedAt time.Time `json:"updatedAt" pg:"updated_at, default:now()"`
}
