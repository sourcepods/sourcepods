package user

type Service interface {
	FindAll() ([]*User, error)
	FindByUsername(string) (*User, error)
	Create(*User) (*User, error)
	Update(string, *User) (*User, error)
	Delete(string) error
}

type Store interface {
	FindAll() ([]*User, error)
	Find(string) (*User, error)
	FindByUsername(string) (*User, error)
	Create(*User) (*User, error)
	Update(string, *User) (*User, error)
	Delete(string) error
}

type service struct {
	users Store
}

func NewService(users Store) Service {
	return &service{users: users}
}

func (s *service) FindAll() ([]*User, error) {
	return s.users.FindAll()
}

func (s *service) FindByUsername(username string) (*User, error) {
	return s.users.FindByUsername(username)
}

func (s *service) Create(user *User) (*User, error) {
	return user, nil
}

func (s *service) Update(username string, user *User) (*User, error) {
	return user, nil
}

func (s *service) Delete(username string) error {
	return nil
}
