package users

import (
	"strconv"
)

// UserJSON ...
type UserJSON struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password,omitempty"`
	Token     string `json:"token,omitempty"`
	CreatedAt int    `json:"created_at"`
}

// UserToJSON ...
func UserToJSON(u *User) *UserJSON {
	obj := NewUserJSON()

	obj.ID = strconv.Itoa(u.ID)
	obj.Email = u.Email
	obj.FirstName = u.FirstName
	obj.LastName = u.LastName
	obj.Token = u.Token
	obj.CreatedAt = u.CreatedAt

	return obj
}

// NewUserJSON ...
func NewUserJSON() *UserJSON {
	json := &UserJSON{}
	return json
}
