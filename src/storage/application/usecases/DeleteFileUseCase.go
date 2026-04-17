package usecases

import (
	"fmt"
	"gestrym-storage/src/storage/domain"
)

type DeleteFileUseCase struct {
	fileRepo       domain.IFileRepository
	storageAdapter domain.IStorageAdapter
}

func NewDeleteFileUseCase(fileRepo domain.IFileRepository, storageAdapter domain.IStorageAdapter) *DeleteFileUseCase {
	return &DeleteFileUseCase{fileRepo: fileRepo, storageAdapter: storageAdapter}
}

func (u *DeleteFileUseCase) Execute(fileID string) error {
	file, err := u.fileRepo.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("could not find file: %v", err)
	}

	// Logic deletion: is_active = false in DB, file remains in storage
	if err := u.fileRepo.Delete(file); err != nil {
		return fmt.Errorf("could not deactivate file metadata: %v", err)
	}

	return nil
}
