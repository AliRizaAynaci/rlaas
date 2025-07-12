package rule

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = gorm.ErrRecordNotFound

type Service struct {
	repo Repository
	db   *gorm.DB // we only need raw DB for owner check
}

func NewService(r Repository, db *gorm.DB) *Service {
	return &Service{repo: r, db: db}
}

/* verifies project.user_id == uid */
func (s *Service) assertOwner(pid, uid uint) error {
	var ownerID uint
	if err := s.db.
		Raw(`SELECT user_id FROM projects WHERE id = ?`, pid).
		Scan(&ownerID).Error; err != nil {
		return err
	}
	if ownerID != uid {
		return errors.New("forbidden")
	}
	return nil
}

/* -------- CRUD wrappers -------- */

func (s *Service) List(uid, pid uint) ([]Rule, error) {
	if err := s.assertOwner(pid, uid); err != nil {
		return nil, err
	}
	return s.repo.ListByProject(pid)
}

func (s *Service) Add(uid, pid uint, in *Rule) (*Rule, error) {
	if err := s.assertOwner(pid, uid); err != nil {
		return nil, err
	}
	in.ProjectID = pid
	return in, s.repo.Create(in)
}

func (s *Service) Update(uid uint, in *Rule) error {
	if err := s.assertOwner(in.ProjectID, uid); err != nil {
		return err
	}
	return s.repo.Update(in)
}

func (s *Service) Delete(uid, pid, rid uint) error {
	if err := s.assertOwner(pid, uid); err != nil {
		return err
	}
	return s.repo.Delete(rid, pid)
}
