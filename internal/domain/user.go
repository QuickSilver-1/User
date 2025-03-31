package domain

import "time"

type User struct {
	Id        uint64     `json:"id"`
	FirstName string     `json:"name"`
	LastName  string     `json:"surname"`
	BirthDay  *time.Time `json:"birthday"`
	Login     string     `json:"email"`
	Password  string     `json:"password"`
}

func NewUser(name, surname, login, pass string, birth *time.Time) *User {
	return &User{
		FirstName: name,
		LastName:  surname,
		BirthDay:  birth,
		Login:     login,
		Password:  pass,
	}
}
