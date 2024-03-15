package models

import (
	"github.com/guregu/null"
	"time"
)

type CaseBase struct {
	CameraID       string    `json:"camera_id" db:"camera_id"`
	Transport      string    `json:"transport" db:"transport"`
	ViolationID    string    `json:"violation_id" db:"violation_id"`
	ViolationValue string    `json:"violation_value" db:"violation_value"`
	Level          int       `json:"level" db:"level"`
	CurrentLevel   int       `json:"current_level" db:"current_level"`
	Datetime       time.Time `json:"datetime" db:"datetime"`
	PhotoUrl       string    `json:"photo_url" db:"photo_url"`
}

type Case struct {
	CaseBase
	ID int `json:"id" db:"id"`
}

type CaseUpdate struct {
	Case
}

type CaseViolations struct {
	Case
	Violation
}

type CaseCursor struct {
	Cases  []CaseViolations `json:"cases"`
	Cursor null.Int         `json:"cursor"`
}

type RatedCreate struct {
	CaseID int  `json:"case_id" db:"case_id" validate:"required"`
	Choice bool `json:"choice" db:"choice"`
}

type RatedUpdate struct {
	CaseID int    `json:"case_id" db:"case_id" validate:"required"`
	Status string `json:"status" db:"status" validate:"required"`
}

type RatedBase struct {
	RatedCreate
	SpecialistID int       `json:"specialist_id" db:"specialist_id"`
	Date         time.Time `json:"date" db:"datetime"`
	Status       string    `json:"status" db:"status"`
}

type Rated struct {
	RatedBase
	Violation
	Level          int    `json:"level"`
	PhotoUrl       string `json:"photo_url"`
	CameraID       string `json:"camera_id"`
	ViolationValue string `json:"violation_value"`
	ID             int    `json:"id"`
}

type RatedCursor struct {
	Rated  []Rated  `json:"rated"`
	Cursor null.Int `json:"cursor"`
}

type RatedCover struct {
	ID     int       `json:"id"`
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
	SpecialistCover
}

type CaseFul struct {
	Violation
	Transport      string       `json:"transport"`
	ViolationValue string       `json:"violation_value"`
	Level          int          `json:"level"`
	Datetime       time.Time    `json:"datetime" `
	PhotoUrl       string       `json:"photo_url"`
	RatedCovers    []RatedCover `json:"rated_covers"`
}
