package repository

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type (
	//ValidationErrors are returned with a slice of all invalid fields
	ValidationErrors struct {
		Errors []ValidationError
	}
	//ValidationError knows for a given field the error
	ValidationError struct {
		Field string
		Error error
	}
)

func (e ValidationErrors) Error() string {
	return fmt.Sprintf("there are %d validation errors", len(e.Errors))
}

// ValidateCreate takes a Repository and validates its fields.
func ValidateCreate(r *Repository) error {
	var errs ValidationErrors

	if err := validateName(r.Name); err != nil {
		errs.Errors = append(errs.Errors, ValidationError{
			Field: "name",
			Error: err,
		})
	}

	if err := validateDescription(r.Description); err != nil {
		errs.Errors = append(errs.Errors, ValidationError{
			Field: "description",
			Error: err,
		})
	}

	if err := validateWebsite(r.Website); err != nil {
		errs.Errors = append(errs.Errors, ValidationError{
			Field: "website",
			Error: err,
		})
	}

	if len(errs.Errors) > 0 {
		return errs
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
