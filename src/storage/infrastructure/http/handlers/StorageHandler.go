package handlers

import (
	"gestrym-storage/src/storage/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	uploadFileUseCase       *usecases.UploadFileUseCase
	getFilesByEntityUseCase *usecases.GetFilesByEntityUseCase
	deleteFileUseCase       *usecases.DeleteFileUseCase
}

func NewStorageHandler(
	uploadFileUseCase *usecases.UploadFileUseCase,
	getFilesByEntityUseCase *usecases.GetFilesByEntityUseCase,
	deleteFileUseCase *usecases.DeleteFileUseCase,
) *StorageHandler {
	return &StorageHandler{
		uploadFileUseCase:       uploadFileUseCase,
		getFilesByEntityUseCase: getFilesByEntityUseCase,
		deleteFileUseCase:       deleteFileUseCase,
	}
}

// UploadFiles godoc
// @Summary Subir uno o varios archivos
// @Description Sube archivos de forma concurrente, valida tipo y tamaño
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Param collection formData string true "Collection/Bucket internal path"
// @Param entityId formData string true "Entity ID to relate (e.g., exercise ID)"
// @Param entityType formData string true "Entity Type to relate (e.g., 'exercise')"
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

	collection := c.PostForm("collection")
	entityID := c.PostForm("entityId")
	entityType := c.PostForm("entityType")

	if collection == "" || entityID == "" || entityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection, entityId, and entityType are required"})
		return
	}

	var requests []usecases.UploadRequest
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file: " + fileHeader.Filename})
			return
		}
		defer file.Close()

		requests = append(requests, usecases.UploadRequest{
			File:       file,
			Header:     fileHeader,
			Collection: collection,
			EntityID:   entityID,
			EntityType: entityType,
		})
	}

	uploadedFiles, err := h.uploadFileUseCase.UploadMultipleFiles(requests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "files uploaded successfully", "data": uploadedFiles})
}

// GetFilesByEntity godoc
// @Summary Obtener archivos por entidad
// @Description Retorna los archivos asociados a una entidad específica con URLs pre-firmadas
// @Tags Storage
// @Produce json
// @Param entityId query string true "Entity ID"
// @Param entityType query string true "Entity Type"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /files [get]
func (h *StorageHandler) GetFilesByEntity(c *gin.Context) {
	entityID := c.Query("entityId")
	entityType := c.Query("entityType")

	if entityID == "" || entityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entityId and entityType are required query parameters"})
		return
	}

	files, err := h.getFilesByEntityUseCase.Execute(entityID, entityType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
