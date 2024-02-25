package swagger

type SpecialistBase struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Fullname string `json:"fullname,omitempty"`
}

type SpecialistCreate struct {
	SpecialistBase
	PhotoUrl string `json:"photoUrl,omitempty"`
}

type Specialist struct {
	SpecialistCreate
	ID         int  `json:"id"`
	Level      int  `json:"level"`
	IsVerified bool `json:"isVerified"`
}

type SpecialistUpdate struct {
	Specialist
}
