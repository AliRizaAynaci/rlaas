package project

type Repository interface {
	Create(*Project) error
	ListByUser(uint) ([]Project, error)
	FindByID(id uint) (*Project, error)
	Delete(id, userID uint) error
}
