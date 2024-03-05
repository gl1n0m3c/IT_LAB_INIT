package models

type Contact struct {
	Transport    string
	UserContacts map[string]string
}
