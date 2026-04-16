package routes

import (
	"gestrym-storage/docs"
	"gestrym-storage/src/common/middleware"
	"gestrym-storage/src/common/utils"
	"net/http"
	"sync"
	"time"

	"gestrym-storage/src/common/config"
	"gestrym-storage/src/storage/application/usecases"
	"gestrym-storage/src/storage/infrastructure/adapters"
	"gestrym-storage/src/storage/infrastructure/http/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type routesDefinition struct {
	serverGroup    *gin.RouterGroup
	publicGroup    *gin.RouterGroup
	privateGroup   *gin.RouterGroup
	internalGroup  *gin.RouterGroup
	protectedGroup *gin.RouterGroup
	logger         utils.ILogger
}

var (
	routesInstance *routesDefinition
	routesOnce     sync.Once
)

func NewRoutesDefinition(serverInstance *gin.Engine) *routesDefinition {
	routesOnce.Do(func() {
		routesInstance = &routesDefinition{}
		routesInstance.logger = utils.NewLogger()
		docs.SwaggerInfo.Title = "Gestrym Training API"
		docs.SwaggerInfo.Description = "API para el manejo de entrenamientos."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.BasePath = "/gestrym-training"
		routesInstance.addCORSConfig(serverInstance)
		routesInstance.addRoutes(serverInstance)
	})
	return routesInstance
}

func (r *routesDefinition) addCORSConfig(serverInstance *gin.Engine) {
	corsMiddleware := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	// Aplica el middleware CORS
	serverInstance.Use(corsMiddleware)
}

func (r *routesDefinition) addRoutes(serverInstance *gin.Engine) {
	// Add default routes
	r.addDefaultRoutes(serverInstance)

	// Instantiate DB
	conn := config.NewPostgresConnection()
	db := conn.GetDB()

	// Repositories
	fileRepo := adapters.NewPostgresFileRepository(db)

	// Adapters & Services
	minioAdapter, err := adapters.NewMinioStorageAdapter()
	if err != nil {
		r.logger.Error("Could not initialize MinIO adapter: " + err.Error())
	}

	uploadFileUseCase := usecases.NewUploadFileUseCase(minioAdapter, fileRepo)
	getFilesUseCase := usecases.NewGetFilesByEntityUseCase(fileRepo, minioAdapter)
	deleteFileUseCase := usecases.NewDeleteFileUseCase(fileRepo, minioAdapter)

	// Controllers
	storageHandler := handlers.NewStorageHandler(uploadFileUseCase, getFilesUseCase, deleteFileUseCase)

	// Add server group
	r.serverGroup = serverInstance.Group(docs.SwaggerInfo.BasePath)
	r.serverGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add groups
	r.publicGroup = r.serverGroup.Group("/public")

	// Register Exercise endpoints
	// (Note: For this service, it is Storage endpoints)
	storageGroup := r.publicGroup.Group("/files")
	{
		storageGroup.POST("/upload", storageHandler.UploadFiles)
		storageGroup.GET("", storageHandler.GetFilesByEntity)
		storageGroup.DELETE("/:id", storageHandler.DeleteFile)
	}

	r.privateGroup = r.serverGroup.Group("/private")
	r.protectedGroup = r.serverGroup.Group("/protected")

	// Add middleware to private group
	r.privateGroup.Use(middleware.SetupJWTMiddleware())

	r.protectedGroup.Use(middleware.SetupApiKeyMiddleware())

	// Add routes to groups
	r.addPublicRoutes()
	r.addPrivateRoutes()
	r.addInternalRoutes()
	r.addProtectedRoutes()

}

func (r *routesDefinition) addDefaultRoutes(serverInstance *gin.Engine) {

	// Handle root
	serverInstance.GET("/", func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "OK",
			"message": "gestrym-training OK...",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusOK, response)
	})

	// Handle 404
	serverInstance.NoRoute(func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "NOT_FOUND",
			"message": "Resource not found",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusNotFound, response)
	})
}

func (r *routesDefinition) addPublicRoutes() {

}

func (r *routesDefinition) addPrivateRoutes() {
}

func (r *routesDefinition) addInternalRoutes() {

}

func (r *routesDefinition) addProtectedRoutes() {
}
