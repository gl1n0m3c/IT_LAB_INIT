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

type RatingSpecialistID struct {
	ID     int
	Level  int
	Rating float32
}

type RatingSpecialistFul struct {
	SpecialistID int
	Level        int
	Fullname     string
	Rating       null.Float
}
