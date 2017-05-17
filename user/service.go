package user

type Service interface {
	FindAll() ([]*User, error)
	FindByUsername(Username) (*User, error)
	Create(*User) (*User, error)
	Update(Username, *User) (*User, error)
	Delete(Username) error
}

type Repository interface {
	FindAll() ([]*User, error)
	Find(ID) (*User, error)
	FindByUsername(Username) (*User, error)
	Create(*User) (*User, error)
	Update(Username, *User) (*User, error)
	Delete(Username) error
}

type service struct {
	users Repository
}

func NewService(users Repository) Service {
	return &service{users: users}
}

func (s *service) FindAll() ([]*User, error) {
	return s.users.FindAll()
}

func (s *service) FindByUsername(username Username) (*User, error) {
	return s.users.FindByUsername(username)
}

func (s *service) Create(user *User) (*User, error) {
	return user, nil
}

func (s *service) Update(username Username, user *User) (*User, error) {
	return user, nil
}

func (s *service) Delete(username Username) error {
	return nil
}
