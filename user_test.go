package gitpods

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	u := &User{}
	err := u.Validate()
	assert.Equal(t, "id: non zero value required;"+
		"username: non zero value required;"+
		"name: non zero value required;"+
		"email: non zero value required;", err.Error())

	// Add values, but not valid ones.
	u.ID = "b755461a-a923-4828-aee1"
	u.Username = "abc"
	u.Name = "bla"
	u.Email = "nomail"
	u.Password = "password"

	err = u.Validate()
	assert.Equal(t, "id: b755461a-a923-4828-aee1 does not validate as uuidv4;"+
		"username: abc does not validate as length(4|32);"+
		"email: nomail does not validate as email;", err.Error())

	// Add valid values
	u.ID = "b755461a-a923-4828-aee1-215903f26e0b"
	u.Username = "metalmatze"
	u.Name = "Matthias Loibl"
	u.Email = "metalmatze@example.com"
	u.Password = "password"
	err = u.Validate()
	assert.NoError(t, err)
}
