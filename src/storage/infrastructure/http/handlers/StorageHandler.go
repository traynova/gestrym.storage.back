package handlers

import (
	"fmt"
	"gestrym-storage/src/common/utils"
	"gestrym-storage/src/storage/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	uploadFileUseCase           *usecases.UploadFileUseCase
	getFilesByCollectionUseCase *usecases.GetFilesByCollectionUseCase
	deleteFileUseCase           *usecases.DeleteFileUseCase
	logger                      utils.ILogger
}

func NewStorageHandler(
	uploadFileUseCase *usecases.UploadFileUseCase,
	getFilesByCollectionUseCase *usecases.GetFilesByCollectionUseCase,
	deleteFileUseCase *usecases.DeleteFileUseCase,
) *StorageHandler {
	return &StorageHandler{
		uploadFileUseCase:           uploadFileUseCase,
		getFilesByCollectionUseCase: getFilesByCollectionUseCase,
		deleteFileUseCase:           deleteFileUseCase,
		logger:                      utils.NewLogger(),
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
// @Security ApiKeyAuth
// @Router /internal/files/upload [post]
func (h *StorageHandler) UploadFiles(c *gin.Context) {
	clientIP := c.ClientIP()
	contentType := c.GetHeader("Content-Type")
	apiKeyPresent := c.GetHeader("X-API-Key") != ""

	h.logger.Info("[UPLOAD] request recibido desde IP:%s | Content-Type:%s | X-API-Key-presente:%v",
		clientIP, contentType, apiKeyPresent)

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error("[UPLOAD] [400] error al parsear multipart form desde IP:%s | Content-Type:%q | error:%v",
			clientIP, contentType, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "could not parse multipart form",
			"detail":       err.Error(),
			"content_type": contentType,
		})
		return
	}

	// Log todos los campos del form para diagnóstico
	for fieldName, fileHeaders := range form.File {
		for _, fh := range fileHeaders {
			h.logger.Info("[UPLOAD] campo file recibido: field=%q filename=%q size=%d contentType=%s",
				fieldName, fh.Filename, fh.Size, fh.Header.Get("Content-Type"))
		}
	}
	for fieldName, values := range form.Value {
		h.logger.Info("[UPLOAD] campo texto recibido: field=%q values=%v", fieldName, values)
	}

	files := form.File["files"]
	if len(files) == 0 {
		availableFields := make([]string, 0, len(form.File))
		for k := range form.File {
			availableFields = append(availableFields, k)
		}
		h.logger.Error("[UPLOAD] [400] campo 'files' no encontrado desde IP:%s | campos file disponibles:%v",
			clientIP, availableFields)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "no files provided",
			"hint":              "El campo multipart debe llamarse 'files' (plural)",
			"campos_recibidos":  fmt.Sprintf("%v", availableFields),
		})
		return
	}

	collectionID := c.PostForm("collectionId")
	service := c.PostForm("service")
	h.logger.Info("[UPLOAD] procesando %d archivo(s) | collectionId:%q service:%q",
		len(files), collectionID, service)

	var requests []usecases.UploadRequest
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			h.logger.Error("[UPLOAD] [500] error al abrir archivo %q: %v", fileHeader.Filename, err)
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
		h.logger.Error("[UPLOAD] [500] error en use case UploadMultipleFiles | collectionId:%q | error:%v",
			resultCollectionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "collection_id": resultCollectionID})
		return
	}

	h.logger.Info("[UPLOAD] [200] archivos subidos exitosamente | collection_id:%q", resultCollectionID)
	c.JSON(http.StatusOK, gin.H{
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
// @Security ApiKeyAuth
// @Router /internal/files/collection [get]
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

	c.JSON(http.StatusOK, gin.H{"data": files[0]})
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
// @Security ApiKeyAuth
// @Router /internal/files/{id} [delete]
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
