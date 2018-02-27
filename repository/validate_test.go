package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreate(t *testing.T) {
	assert.Error(t, ValidateCreate(&Repository{}))
	assert.Error(t, ValidateCreate(&Repository{
		Name: "foo",
	}))
	assert.Error(t, ValidateCreate(&Repository{
		Name:    "foobar",
		Website: "bar",
	}))

	assert.Nil(t, ValidateCreate(&Repository{
		Name: "foobar",
	}))
	assert.Nil(t, ValidateCreate(&Repository{
		Name:        "foobar",
		Description: "example",
		Website:     "http://example.com",
	}))
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
