package models

import "github.com/guregu/null"

type RatingSpecialistCount struct {
	SpecialistCover
	Total   int `json:"total"`
	Correct int `json:"correct"`
	Unknown int `json:"unknown"`
}

type RatingSpecialistCountCursor struct {
	Specialists []RatingSpecialistCount `json:"specialists"`
	Cursor      null.Int                `json:"cursor"`
}
