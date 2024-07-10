package model

type User struct {
	Id             int    `json:"id" example:"1"`
	PassportNumber string `json:"passportNumber" example:"1234 567890"`
	Name           string `json:"name" example:"Petr"`
	Surname        string `json:"surname" example:"Petr"`
	Patronymic     string `json:"patronymic,omitempty" example:"Petr"`
	Address        string `json:"address" example:"Piter"`
} // @name User
