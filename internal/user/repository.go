package user

type Repository interface {
	Create(u *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(e string) (*User, error)
	FindByGoogleID(gid string) (*User, error)
	UpdateNameAndPic(id uint, name, pic string) error
}
