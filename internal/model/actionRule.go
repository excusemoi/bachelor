package model

import "time"

type ActionRule struct {
	tableName        struct{}  `json:"-" pg:"action_rule,discard_unknown_columns" gorm:"autoIncrement"`
	Id               int       `json:"id" pg:"id,pk"`
	BootstrapServers string    `json:"bootstrapServers" pg:"bootstrap_servers"`
	Key              string    `json:"key" pg:"key"`
	UpdatedAt        time.Time `json:"updatedAt" pg:"updated_at, default:now()"`
}
