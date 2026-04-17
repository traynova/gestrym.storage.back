package usecases

import (
	"fmt"
	"gestrym-storage/src/common/models"
	"gestrym-storage/src/storage/domain"
)

type GetFilesByCollectionUseCase struct {
	fileRepo       domain.IFileRepository
	storageAdapter domain.IStorageAdapter
}

func NewGetFilesByCollectionUseCase(fileRepo domain.IFileRepository, storageAdapter domain.IStorageAdapter) *GetFilesByCollectionUseCase {
	return &GetFilesByCollectionUseCase{fileRepo: fileRepo, storageAdapter: storageAdapter}
}

func (u *GetFilesByCollectionUseCase) Execute(collectionID string) ([]models.File, error) {
	files, err := u.fileRepo.FindByCollectionID(collectionID)
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
