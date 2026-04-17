package domain

import "gestrym-storage/src/common/models"

type IFileRepository interface {
	Save(file *models.File) error
	FindByID(id string) (*models.File, error)
	FindByEntity(entityID, entityType string) ([]models.File, error)
	FindByCollectionID(collectionID string) ([]models.File, error)
	Delete(file *models.File) error
}
