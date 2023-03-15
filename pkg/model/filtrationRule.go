package model

type FiltrationRule struct {
	FilterId       int    `json:"filterId" pg:"filter_id"`
	RuleId         int    `json:"ruleId" pg:"rule_id"`
	FilterField    string `json:"filterField" pg:"filter_field"`
	FilterFunction string `json:"filterFunction" pg:"filter_function"`
	FilterValue    string `json:"filterValue" pg:"filter_value"`
}
