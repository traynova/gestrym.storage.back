package domain

import "io"

type IStorageAdapter interface {
	UploadFile(fileStream io.Reader, objectSize int64, contentType, fileName, collection string) (string, error)
	DeleteFile(fileName, collection string) error
	GetFileURL(fileName, collection string) (string, error)
}
