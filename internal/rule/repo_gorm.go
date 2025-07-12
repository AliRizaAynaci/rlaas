package rule

import "gorm.io/gorm"

type gormRepo struct{ db *gorm.DB }

func NewGormRepo(db *gorm.DB) Repository { return &gormRepo{db} }

func (r *gormRepo) Create(m *Rule) error { return r.db.Create(m).Error }

func (r *gormRepo) ListByProject(pid uint) ([]Rule, error) {
	var rs []Rule
	return rs, r.db.Where("project_id=?", pid).Find(&rs).Error
}

func (r *gormRepo) Update(m *Rule) error {
	return r.db.Where("id=? AND project_id=?", m.ID, m.ProjectID).Updates(m).Error
}

func (r *gormRepo) Delete(id, pid uint) error {
	return r.db.Where("id=? AND project_id=?", id, pid).Delete(&Rule{}).Error
}
