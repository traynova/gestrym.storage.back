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

	if err := u.storageAdapter.DeleteFile(file.FileName, file.Collection); err != nil {
		return fmt.Errorf("could not delete file from storage: %v", err)
	}

	if err := u.fileRepo.Delete(file); err != nil {
		return fmt.Errorf("could not delete file metadata: %v", err)
	}

	return nil
}
