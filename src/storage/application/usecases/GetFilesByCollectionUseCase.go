package usecases

import (
	"fmt"
	"gestrym-storage/src/storage/domain"
	"gestrym-storage/src/storage/infrastructure/structs"
)

type GetFilesByCollectionUseCase struct {
	fileRepo       domain.IFileRepository
	storageAdapter domain.IStorageAdapter
}

func NewGetFilesByCollectionUseCase(fileRepo domain.IFileRepository, storageAdapter domain.IStorageAdapter) *GetFilesByCollectionUseCase {
	return &GetFilesByCollectionUseCase{fileRepo: fileRepo, storageAdapter: storageAdapter}
}

func (u *GetFilesByCollectionUseCase) Execute(collectionID string) ([]structs.GetResponse, error) {
	files, err := u.fileRepo.FindByCollectionID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve files from repository: %v", err)
	}

	var response []structs.GetResponse
	for _, file := range files {
		url, err := u.storageAdapter.GetFileURL(file.FileName, file.CollectionID)
		if err == nil {
			response = append(response, structs.GetResponse{
				Id:           file.ID,
				FileName:     file.FileName,
				ContentType:  file.ContentType,
				Size:         file.Size,
				URL:          url,
				CollectionID: file.CollectionID,
			})
		}
	}

	return response, nil
}
