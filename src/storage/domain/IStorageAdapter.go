package domain

import "io"

type IStorageAdapter interface {
	UploadFile(fileStream io.Reader, objectSize int64, contentType, fileName, collectionID string) (string, error)
	DeleteFile(fileName, collectionID string) error
	GetFileURL(fileName, collectionID string) (string, error)
}
