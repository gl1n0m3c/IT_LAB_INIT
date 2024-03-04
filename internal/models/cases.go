package models

import "time"

type CaseBase struct {
	CameraID       string    `json:"camera_id" db:"camera_id"`
	Transport      string    `json:"transport" db:"transport"`
	ViolationID    string    `json:"violation_id" db:"violation_id"`
	ViolationValue string    `json:"violation_value" db:"violation_value"`
	Level          int       `json:"level" db:"level"`
	Datetime       time.Time `json:"datetime" db:"datetime"`
	PhotoUrl       string    `json:"photo_url" db:"photo_url"`
}

type Case struct {
	CaseBase
	ID int
}
