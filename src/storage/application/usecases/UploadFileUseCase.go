package usecases

import (
	"fmt"
	"gestrym-storage/src/common/models"
	"gestrym-storage/src/storage/domain"
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
	File        multipart.File
	Header      *multipart.FileHeader
	Collection  string
	EntityID    string
	EntityType  string
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

	fileName := fmt.Sprintf("%s-%s", req.EntityID, req.Header.Filename)
	
	objectName, err := u.storageAdapter.UploadFile(req.File, req.Header.Size, contentType, fileName, req.Collection)
	if err != nil {
		return nil, fmt.Errorf("could not upload file: %v", err)
	}

	_, err = u.storageAdapter.GetFileURL(fileName, req.Collection)
	if err != nil {
		return nil, fmt.Errorf("could not get file url: %v", err)
	}

	// Remove query params from URL for storage if it's public, or just store the object path
	// Assuming MinIO presigned URL is fine for now, or you could store just object Name and generate URL on the fly.
	// We'll store the objectName for URL so it can be generated or accessed consistently.
	
	newFile := &models.File{
		FileName:    fileName,
		ContentType: contentType,
		Size:        req.Header.Size,
		URL:         objectName, // Store path, generate presigned url on access
		Collection:  req.Collection,
		EntityID:    req.EntityID,
		EntityType:  req.EntityType,
	}

	if err := u.fileRepo.Save(newFile); err != nil {
		u.storageAdapter.DeleteFile(fileName, req.Collection) // Rollback
		return nil, fmt.Errorf("could not save file metadata: %v", err)
	}

	return newFile, nil
}

// UploadMultipleFiles handles concurrent uploading of multiple files
func (u *UploadFileUseCase) UploadMultipleFiles(requests []UploadRequest) ([]*models.File, error) {
	var wg sync.WaitGroup
	results := make([]*models.File, len(requests))
	errorsChan := make(chan error, len(requests))

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request UploadRequest) {
			defer wg.Done()
			file, err := u.UploadSingleFile(request)
			if err != nil {
				errorsChan <- err
				return
			}
			results[index] = file
		}(i, req)
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
		return nil, fmt.Errorf("encountered errors during concurrent upload: %v", uploadErrors)
	}

	return results, nil
}
