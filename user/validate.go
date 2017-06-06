package user

import (
	"fmt"

	"gopkg.in/asaskevich/govalidator.v6"
)

func ValidateCreate(u *User) []error {
	var errs []error

	if err := validateEmail(u.Email); err != nil {
		errs = append(errs, err)
	}

	if err := validateUsername(u.Username); err != nil {
		errs = append(errs, err)
	}

	if err := validateName(u.Name); err != nil {
		errs = append(errs, err)
	}

	if err := validatePassword(u.Password); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func validateID(id string) error {
	if ok := govalidator.IsUUIDv4(id); !ok {
		return fmt.Errorf("id is not a valid uuid v4")
	}
	return nil
}

func validateEmail(email string) error {
	if ok := govalidator.IsEmail(email); !ok {
		return fmt.Errorf("email is not valid")
	}
	return nil
}

func validateUsername(username string) error {
	if ok := govalidator.IsAlphanumeric(username); !ok {
		return fmt.Errorf("username is not alphanumeric")
	}
	if ok := govalidator.IsByteLength(username, 4, 32); !ok {
		return fmt.Errorf("username is not between 4 and 32 characters long")
	}
	return nil
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("name can't be empty")
	}
	return nil
}

func validatePassword(pass string) error {
	if len(pass) < 6 {
		return fmt.Errorf("password needs to be at least 6 characters")
	}
	return nil
}
