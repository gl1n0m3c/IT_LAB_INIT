package models

type Violation struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}
