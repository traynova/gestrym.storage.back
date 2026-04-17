package handlers

import (
	"gestrym-storage/src/storage/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	uploadFileUseCase            *usecases.UploadFileUseCase
	getFilesByCollectionUseCase *usecases.GetFilesByCollectionUseCase
	deleteFileUseCase            *usecases.DeleteFileUseCase
}

func NewStorageHandler(
	uploadFileUseCase *usecases.UploadFileUseCase,
	getFilesByCollectionUseCase *usecases.GetFilesByCollectionUseCase,
	deleteFileUseCase *usecases.DeleteFileUseCase,
) *StorageHandler {
	return &StorageHandler{
		uploadFileUseCase:            uploadFileUseCase,
		getFilesByCollectionUseCase: getFilesByCollectionUseCase,
		deleteFileUseCase:            deleteFileUseCase,
	}
}

// UploadFiles godoc
// @Summary Subir uno o varios archivos
// @Description Sube archivos de forma concurrente, valida tipo y tamaño
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Param collectionId formData string false "Collection ID to relate (optional, will generate one if empty)"
// @Param service formData string false "Service of origin (e.g., 'users', 'exercises')"
// @Param files formData file true "Archivos a subir"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /files/upload [post]
func (h *StorageHandler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse multipart form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
		return
	}

	collectionID := c.PostForm("collectionId")
	service := c.PostForm("service")

	var requests []usecases.UploadRequest
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file: " + fileHeader.Filename})
			return
		}
		defer file.Close()

		requests = append(requests, usecases.UploadRequest{
			File:         file,
			Header:       fileHeader,
			CollectionID: collectionID,
			Service:      service,
		})
	}

	resultCollectionID, err := h.uploadFileUseCase.UploadMultipleFiles(requests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "collection_id": resultCollectionID})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "files uploaded successfully",
		"collection_id": resultCollectionID,
	})
}


// GetFilesByCollection godoc
// @Summary Obtener archivos por ID de colección
// @Description Retorna los archivos asociados a una colección específica
// @Tags Storage
// @Produce json
// @Param collectionId query string true "Collection ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /files/collection [get]
func (h *StorageHandler) GetFilesByCollection(c *gin.Context) {
	collectionID := c.Query("collectionId")

	if collectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectionId is required"})
		return
	}

	files, err := h.getFilesByCollectionUseCase.Execute(collectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If only one file, we can still return as a list or just the object,
	// but keeping consistent with a list is usually better.
	// The requirement: "si es solo uno se retorna ese" might mean don't return an array if size is 1.
	if len(files) == 1 {
		c.JSON(http.StatusOK, gin.H{"data": files[0]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}

// DeleteFile godoc
// @Summary Eliminar archivo
// @Description Elimina un archivo del storage y de la base de datos
// @Tags Storage
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /files/{id} [delete]
func (h *StorageHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")

	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file ID is required"})
		return
	}

	if err := h.deleteFileUseCase.Execute(fileID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}
