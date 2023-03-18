package model

import "time"

type FiltrationRule struct {
	tableName      struct{}  `json:"-" pg:"filtration_rule,discard_unknown_columns"`
	Id             int       `json:"id" pg:",pk"`
	Filter         *Filter   `json:"filter" pg:"-"`
	FilterId       int       `json:"-" pg:"filter_id"`
	Rule           *Rule     `json:"rule" pg:"-"`
	RuleId         int       `json:"-" pg:"rule_id"`
	FilterField    string    `json:"filterField" pg:"filter_field"`
	FilterFunction string    `json:"filterFunction" pg:"filter_function"`
	FilterValue    string    `json:"filterValue" pg:"filter_value"`
	UpdatedAt      time.Time `json:"updatedAt" pg:""`
}
