package gitpods

// User of GitPods.
type User struct {
	ID       string `json:"id"`       // valid:"required,uuidv4"
	Email    string `json:"email"`    // valid:"required,email"
	Username string `json:"username"` // valid:"required,alphanum,length(4|32)"
	Name     string `json:"name"`     // valid:"required"
	Password string `json:"-"`

	Repositories []*Repository `json:"repositories,omitempty"`
}
