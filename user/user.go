package user

type ID string
type Username string

type User struct {
	ID       ID       `json:"id"`       // valid:"required,uuidv4"
	Username Username `json:"username"` // valid:"required,alphanum,length(4|32)"
	Name     string   `json:"name"`     // valid:"required"
	Email    string   `json:"email"`    // valid:"required,email"
	Password string   `json:"-"`

	//Repositories []*Repository `json:"repositories,omitempty"`
}
