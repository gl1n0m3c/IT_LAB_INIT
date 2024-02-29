package models

type CameraBase struct {
	Type        string     `json:"type" db:"type" validate:"required"`
	Coordinates [2]float64 `json:"coordinates" db:"coordinates" validate:"required"`
	Description string     `json:"description" db:"description" validate:"required"`
}

type Camera struct {
	ID int `json:"id" db:"id"`
	CameraBase
}
