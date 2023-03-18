package model

type Filter struct {
	tableName struct{} `json:"-" pg:"filter,discard_unknown_columns"`
	Id        int      `json:"id" pg:",pk"`
	Filter    string   `json:"filter" pg:"filter"`
}
