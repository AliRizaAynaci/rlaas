package rule

type Repository interface {
	Create(*Rule) error
	ListByProject(uint) ([]Rule, error)
	Update(*Rule) error
	Delete(ruleID, projectID uint) error
}
