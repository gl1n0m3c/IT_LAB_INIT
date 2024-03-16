package models

import (
	"github.com/guregu/null"
)

type SpecialistBase struct {
	Login    string      `db:"login" json:"login" validate:"required"`
	Password string      `db:"hashed_password" json:"password" validate:"required,password"`
	Fullname null.String `db:"fullname" json:"fullname,omitempty"`
}

type SpecialistCreate struct {
	SpecialistBase
	PhotoUrl null.String `db:"photo_url" json:"photoUrl,omitempty"`
}

type Specialist struct {
	SpecialistCreate
	ID         int  `db:"id" json:"id"`
	Level      int  `db:"level" json:"level"`
	Row        int  `db:"row" json:"row"`
	IsVerified bool `db:"is_verified" json:"isVerified"`
}

type SpecialistUpdate struct {
	ID       int    `json:"id" db:"id"`
	Password string `json:"password" db:"hashed_password" validate:"password"`
	FullName string `json:"full_name" db:"fullname"`
	PhotoUrl string `json:"photo_url" db:"photo_url"`
}

type SpecialistLogin struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SpecialistCover struct {
	ID       int         `json:"id"`
	Fullname null.String `json:"fullname"`
	Level    int         `json:"level"`
	Row      int         `json:"row"`
	PhotoUrl null.String `json:"photo_url"`
}
