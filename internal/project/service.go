package project

type Service struct{ repo Repository }

func NewService(r Repository) *Service { return &Service{r} }

func (s *Service) Create(userID uint, name string, apikey string) (*Project, error) {
	p := &Project{
		UserID: userID,
		Name:   name,
		APIKey: apikey,
	}
	return p, s.repo.Create(p)
}

func (s *Service) List(userID uint) ([]Project, error) {
	return s.repo.ListByUser(userID)
}

func (s *Service) UserOwns(pid, uid uint) (bool, error) {
	p, err := s.repo.FindByID(pid)
	if err != nil {
		return false, err
	}
	return p.UserID == uid, nil
}

func (s *Service) Delete(userID, projectID uint) error {
	return s.repo.Delete(projectID, userID)
}
