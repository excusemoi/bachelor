package filter

import (
	"sync"
	"time"
)

type FiltrationRule struct {
	tableName      struct{}     `json:"-" pg:"filtration_rule,discard_unknown_columns"`
	Id             int          `json:"id" pg:",pk"`
	Filter         *Filter      `json:"filter" pg:"-"`
	FilterId       int          `json:"-" pg:"filter_id"`
	FilterField    string       `json:"filterField" pg:"filter_field"`
	FilterFunction string       `json:"filterFunction" pg:"filter_function"`
	FilterValue    string       `json:"filterValue" pg:"filter_value"`
	UpdatedAt      time.Time    `json:"updatedAt" pg:"updated_at"`
	Mx             sync.RWMutex `json:"-" pg:"-"`
}