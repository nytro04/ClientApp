package users

//features of a user s
type User struct {
	ID   int64
	Name string
	Hash []byte
}

//User interface to abstract functionality
type UserDB interface {
	CreateUser(*User) (int64, error)
	GetUserByID(id int64) (*User, error)
	GetUserByName(name string) (*User, error)
	GetAllUsers() ([]*User, error)
	UpdateUser(*User) error
	RemoveUser(id int64) error
}
