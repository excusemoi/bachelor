package model

import (
	"errors"
	"strings"
	"time"
)

const (
	equals      = "equals"
	notEquals   = "notEquals"
	contains    = "contains"
	notContains = "notContains"
)

type FiltrationRule struct {
	tableName struct{}  `json:"-" pg:"filtration_rule,discard_unknown_columns"`
	Id        int       `json:"id" pg:"id,pk"`
	Field     string    `json:"field" pg:"field"`
	Function  string    `json:"function" pg:"function"`
	Value     string    `json:"value" pg:"value"`
	UpdatedAt time.Time `json:"updatedAt" pg:"updated_at, default:now()"`
}

func (r *FiltrationRule) FilterFunction(value string) (bool, error) {
	switch r.Function {
	case equals:
		return r.equals(value), nil
	case notEquals:
		return !r.equals(value), nil
	case contains:
		return r.contains(value), nil
	case notContains:
		return !r.contains(value), nil
	}
	return false, errors.New("incorrect function name")
}

func (r *FiltrationRule) equals(value string) bool {
	return r.Value == value
}

func (r *FiltrationRule) contains(value string) bool {
	return strings.Contains(value, r.Value)
}
