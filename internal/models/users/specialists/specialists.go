package specialists

type specialistBase struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

type SpecialistCreate struct {
	specialistBase
	PhotoUrl string `json:"photoUrl"`
}

type Specialist struct {
	SpecialistCreate
	ID         int  `json:"id"`
	Level      int  `json:"level"`
	IsVerified bool `json:"isVerified"`
}

type SpecialistUpdate struct {
	SpecialistCreate
	ID int `json:"id"`
}
