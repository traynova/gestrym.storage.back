package config

import (
	"fmt"
	"gestrym-storage/src/common/models"
	"gestrym-storage/src/common/utils"
)

var logger = utils.NewLogger()

func MigrateDB() (IDatabaseConnection, error) {
	connection := NewPostgresConnection()
	db := connection.GetDB()

	//Se agregan los modelos de base de datos
	err := db.AutoMigrate()

	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] Error al migrar las entidades: %s", err.Error()))
		return nil, err
	}

	// Agrega un indice único a la tabla 'permissions'
	// haciendo que la combinación de 'role', 'action' y 'resource' sea única
	// para evitar convinaciones duplicadas
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_role_action_resource ON permissions (role_id, action_id, resource_id)").Error
	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] Error al agregar el índice único a la tabla de permisos: %s", err.Error()))
		return nil, err
	}

	defaultTokenTypes := []models.UserTokenType{
		{Type: models.UserTokenTypePasswordRecovery},
		{Type: models.UserTokenTypeActivation},
	}
	for _, tokenType := range defaultTokenTypes {
		result := &models.UserTokenType{}
		if err := db.Where("type = ?", tokenType.Type).FirstOrCreate(result, models.UserTokenType{Type: tokenType.Type}).Error; err != nil {
			logger.Error(fmt.Sprintf("[ERROR] Error al insertar los tipos de token de usuario: %s", err.Error()))
			return nil, err
		}
	}

	logger.Info("[OK] Todas las migraciones completadas exitosamente")
	return connection, nil
}
