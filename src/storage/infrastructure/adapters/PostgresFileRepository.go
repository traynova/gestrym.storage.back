package adapters

import (
	"gestrym-storage/src/common/models"
	"gestrym-storage/src/storage/domain"

	"gorm.io/gorm"
)

type postgresFileRepository struct {
	db *gorm.DB
}

func NewPostgresFileRepository(db *gorm.DB) domain.IFileRepository {
	return &postgresFileRepository{db: db}
}

func (r *postgresFileRepository) Save(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *postgresFileRepository) FindByID(id string) (*models.File, error) {
	var file models.File
	if err := r.db.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *postgresFileRepository) FindByEntity(entityID, entityType string) ([]models.File, error) {
	var files []models.File
	if err := r.db.Where("entity_id = ? AND entity_type = ?", entityID, entityType).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *postgresFileRepository) Delete(file *models.File) error {
	return r.db.Delete(file).Error
}
