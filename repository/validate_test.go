package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreate(t *testing.T) {
	testcases := []struct {
		Name       string
		Repository *Repository
		Errors     []ValidationError
	}{
		{
			Name:       "NoInput",
			Repository: &Repository{},
			Errors: []ValidationError{{
				Field: "name",
				Error: errors.New("name is not between 4 and 32 characters long"),
			}},
		},
		{
			Name:       "NameTooShort",
			Repository: &Repository{Name: "foo"},
			Errors: []ValidationError{{
				Field: "name",
				Error: errors.New("name is not between 4 and 32 characters long"),
			}},
		},
		{
			Name:       "NameTooLong",
			Repository: &Repository{Name: "thisnameiswaytolongtobeadecentusername"},
			Errors: []ValidationError{{
				Field: "name",
				Error: errors.New("name is not between 4 and 32 characters long"),
			}},
		},
		{
			Name:       "NameTooLong",
			Repository: &Repository{Name: "thisnameiswaytolongtobeadecentusername"},
			Errors: []ValidationError{{
				Field: "name",
				Error: errors.New("name is not between 4 and 32 characters long"),
			}},
		},
		{
			Name:       "InvalidWebsite",
			Repository: &Repository{Name: "username", Website: "example"},
			Errors: []ValidationError{{
				Field: "website",
				Error: errors.New("example is not a url"),
			}},
		},
		{
			Name: "Valid",
			Repository: &Repository{
				Name:        "username",
				Website:     "http://example.com",
				Description: "Awesome repository!",
			},
			Errors: nil,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			err := ValidateCreate(tc.Repository)

			if len(tc.Errors) > 0 {
				assert.Error(t, err)
				assert.Equal(t, fmt.Sprintf("there are %d validation errors", len(tc.Errors)), err.Error())
				assert.IsType(t, ValidationErrors{}, err)
				if verr, ok := err.(ValidationErrors); ok {
					assert.Len(t, verr.Errors, len(tc.Errors))
					for i, expected := range tc.Errors {
						assert.Equal(t, expected, verr.Errors[i])
					}
				}
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestValidateName(t *testing.T) {
	assert.Error(t, validateName(""))
	assert.Error(t, validateName("a"))
	assert.Error(t, validateName("aa"))
	assert.Error(t, validateName("aaa"))
	assert.Nil(t, validateName("aaaa"))
	assert.Nil(t, validateName("aaaaasdf"))
	assert.Error(t, validateName("aaaasdf!@#~"))
}

func TestValidateDescription(t *testing.T) {
	assert.Nil(t, validateDescription(""))
	assert.Nil(t, validateDescription("asdf"))
	assert.Nil(t, validateDescription("asdf"))
	assert.Nil(t, validateDescription("asdf asd fasdf asdf "))
}

func TestValidateWebsite(t *testing.T) {
	// website is optional and thus the string can be empty
	assert.Nil(t, validateWebsite(""))

	assert.Error(t, validateWebsite("asdf"))
	assert.Error(t, validateWebsite("http://"))

	assert.Nil(t, validateWebsite("http://localhost"))
	assert.Nil(t, validateWebsite("http://example.com"))
	assert.Nil(t, validateWebsite("https://example.com"))
}
