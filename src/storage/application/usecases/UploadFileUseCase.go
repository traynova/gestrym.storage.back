package usecases

import (
	"fmt"
	"gestrym-storage/src/common/models"
	"gestrym-storage/src/storage/domain"
	"github.com/google/uuid"
	"mime/multipart"
	"sync"
)

const MaxFileSize = 10 * 1024 * 1024 // 10MB

var AllowedContentTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

type UploadFileUseCase struct {
	storageAdapter domain.IStorageAdapter
	fileRepo       domain.IFileRepository
}

func NewUploadFileUseCase(storageAdapter domain.IStorageAdapter, fileRepo domain.IFileRepository) *UploadFileUseCase {
	return &UploadFileUseCase{
		storageAdapter: storageAdapter,
		fileRepo:       fileRepo,
	}
}

type UploadRequest struct {
	File         multipart.File
	Header       *multipart.FileHeader
	Collection   string
	CollectionID string
	EntityID     string
	EntityType   string
}

// UploadSingleFile handles validation and uploading of a single file
func (u *UploadFileUseCase) UploadSingleFile(req UploadRequest) (*models.File, error) {
	if req.Header.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds the maximum limit of 10MB")
	}

	contentType := req.Header.Header.Get("Content-Type")
	if !AllowedContentTypes[contentType] {
		return nil, fmt.Errorf("invalid file type: %s", contentType)
	}

	// Generate unique filename to avoid collisions in same collection/bucket
	fileName := fmt.Sprintf("%s-%s", uuid.New().String(), req.Header.Filename)

	objectName, err := u.storageAdapter.UploadFile(req.File, req.Header.Size, contentType, fileName, req.Collection)
	if err != nil {
		return nil, fmt.Errorf("could not upload file: %v", err)
	}

	newFile := &models.File{
		FileName:     fileName,
		ContentType:  contentType,
		Size:         req.Header.Size,
		URL:          objectName, // Store path/object name
		Collection:   req.Collection,
		CollectionID: req.CollectionID,
		EntityID:     req.EntityID,
		EntityType:   req.EntityType,
		IsActive:     true,
	}

	if err := u.fileRepo.Save(newFile); err != nil {
		u.storageAdapter.DeleteFile(fileName, req.Collection) // Rollback
		return nil, fmt.Errorf("could not save file metadata: %v", err)
	}

	return newFile, nil
}

// UploadMultipleFiles handles concurrent uploading of multiple files and returns the collection ID
func (u *UploadFileUseCase) UploadMultipleFiles(requests []UploadRequest) (string, error) {
	if len(requests) == 0 {
		return "", fmt.Errorf("no files to upload")
	}

	// All files in this batch will share the same CollectionID if not provided
	collectionID := requests[0].CollectionID
	if collectionID == "" {
		collectionID = uuid.New().String()
	}

	var wg sync.WaitGroup
	errorsChan := make(chan error, len(requests))

	for i := range requests {
		requests[i].CollectionID = collectionID // Ensure all share the same ID
		wg.Add(1)
		go func(req UploadRequest) {
			defer wg.Done()
			_, err := u.UploadSingleFile(req)
			if err != nil {
				errorsChan <- err
			}
		}(requests[i])
	}

	wg.Wait()
	close(errorsChan)

	var uploadErrors []error
	for err := range errorsChan {
		if err != nil {
			uploadErrors = append(uploadErrors, err)
		}
	}

	if len(uploadErrors) > 0 {
		return collectionID, fmt.Errorf("encountered errors during concurrent upload: %v", uploadErrors)
	}

	return collectionID, nil
}
