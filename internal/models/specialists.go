package models

import (
	"github.com/guregu/null"
)

type SpecialistBase struct {
	Login    string      `db:"login" json:"login"`
	Password string      `db:"hashed_password" json:"password"`
	Fullname null.String `db:"fullname" json:"fullname"`
}

type SpecialistCreate struct {
	SpecialistBase
	PhotoUrl null.String `db:"photo_url" json:"photoUrl"`
}

type Specialist struct {
	SpecialistCreate
	ID         int  `db:"id" json:"id"`
	Level      int  `db:"level" json:"level"`
	IsVerified bool `db:"is_verified" json:"isVerified"`
}

type SpecialistUpdate struct {
	Specialist
}
