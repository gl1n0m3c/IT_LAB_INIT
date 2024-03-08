package models

type ManagerBase struct {
	Login    string `json:"login" db:"login" validate:"required"`
	Password string `json:"password" db:"hashed_password" validate:"required"`
}

type Manager struct {
	ManagerBase
	ID int `json:"id" db:"id"`
}
