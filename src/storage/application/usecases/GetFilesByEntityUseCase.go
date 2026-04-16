package usecases

import (
	"fmt"
	"gestrym-storage/src/common/models"
	"gestrym-storage/src/storage/domain"
)

type GetFilesByEntityUseCase struct {
	fileRepo       domain.IFileRepository
	storageAdapter domain.IStorageAdapter
}

func NewGetFilesByEntityUseCase(fileRepo domain.IFileRepository, storageAdapter domain.IStorageAdapter) *GetFilesByEntityUseCase {
	return &GetFilesByEntityUseCase{fileRepo: fileRepo, storageAdapter: storageAdapter}
}

func (u *GetFilesByEntityUseCase) Execute(entityID, entityType string) ([]models.File, error) {
	files, err := u.fileRepo.FindByEntity(entityID, entityType)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve files from repository: %v", err)
	}

	for i, file := range files {
		url, err := u.storageAdapter.GetFileURL(file.FileName, file.Collection)
		if err == nil {
			files[i].URL = url
		}
	}

	return files, nil
}
