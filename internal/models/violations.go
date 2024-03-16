package models

import "time"

type Violation struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

type FineData struct {
	Violation
	Mail           string
	PhotoUrl       string
	Coordinated    string
	ViolationValue string
	Date           time.Time
}
