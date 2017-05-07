package gitpods

import (
	"gopkg.in/asaskevich/govalidator.v6"
)

type User struct {
	ID       string `json:"id" valid:"required,uuidv4"`
	Username string `json:"username" valid:"required,alphanum,length(4|32)"`
	Name     string `json:"name" valid:"required"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"-"`
}

func (u *User) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	return err
}
