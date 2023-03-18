package model

type Rule struct {
	tableName struct{} `json:"-" pg:"rule,discard_unknown_columns"`
	Id        int      `json:"id" pg:",pk"`
	Rule      string   `json:"rule" pg:"rule"`
}
