package repository

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

func ValidateCreate(r *Repository) error {
	if err := validateName(r.Name); err != nil {
		return err
	}

	if err := validateDescription(r.Description); err != nil {
		return err
	}

	if err := validateWebsite(r.Website); err != nil {
		return err
	}

	return nil
}

func validateID(id string) error {
	if ok := govalidator.IsUUIDv4(id); !ok {
		return fmt.Errorf("id is not a valid uuid v4")
	}
	return nil
}

func validateName(name string) error {
	if ok := govalidator.IsAlphanumeric(name); !ok {
		return fmt.Errorf("name is not alphanumeric")
	}
	if ok := govalidator.IsByteLength(name, 4, 32); !ok {
		return fmt.Errorf("name is not between 4 and 32 characters long")
	}
	return nil
}

func validateDescription(description string) error {
	return nil
}

func validateWebsite(website string) error {
	// website is optional and thus can be an empty string
	if website == "" {
		return nil
	}

	if ok := govalidator.IsURL(website); !ok {
		return fmt.Errorf("%s is not a url", website)
	}

	return nil
}
