package model

type User struct {
	Id             int    `json:"id" example:"1"`
	PassportNumber string `json:"passportNumber" example:"1234 567890"`
	Name           string `json:"name" example:"string"`
	Surname        string `json:"surname" example:"string"`
	Patronymic     string `json:"patronymic,omitempty" example:"string"`
	Address        string `json:"address" example:"string"`
} // @name User
