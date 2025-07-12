package project

import "gorm.io/gorm"

type gormRepo struct{ db *gorm.DB }

func NewGormRepo(db *gorm.DB) Repository { return &gormRepo{db} }

func (r *gormRepo) Create(p *Project) error {
	return r.db.Create(p).Error
}

func (r *gormRepo) ListByUser(uid uint) ([]Project, error) {
	var list []Project
	return list, r.db.
		Preload("Rules").
		Where("user_id = ?", uid).
		Order("created_at DESC").
		Find(&list).Error
}

func (r *gormRepo) FindByID(id uint) (*Project, error) {
	var p Project
	err := r.db.First(&p, id).Error
	return &p, err
}

func (r *gormRepo) Delete(id, uid uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, uid).Delete(&Project{}).Error
}
