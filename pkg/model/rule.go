package model

import "sync"

type Rule struct {
	tableName struct{}     `json:"-" pg:"rule,discard_unknown_columns"`
	Id        int          `json:"id" pg:",pk"`
	Rule      string       `json:"rule" pg:"rule"`
	Mx        sync.RWMutex `json:"-" pg:"-"`
}
