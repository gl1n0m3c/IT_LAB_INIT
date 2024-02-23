package managers

type managerBase struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ManagerCreate struct {
	managerBase
}

type Manager struct {
	managerBase
	ID int `json:"id"`
}

type ManagerUpdate struct {
	ManagerCreate
	ID int `json:"id"`
}
