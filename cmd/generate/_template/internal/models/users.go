package models

type Users struct {
	ID        string `json:"id"`
	CHI       string `json:"chi"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
