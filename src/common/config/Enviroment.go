package config

import (
	"errors"
	"gestrym-storage/src/common/utils"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type env struct {
	GIN_MODE               string `mapstructure:"GIN_MODE" validate:"required,oneof=debug release test"`
	GOOGLE_CLIENT_ID       string `mapstructure:"GOOGLE_CLIENT_ID"`
	GORM_LOG_LEVEL         string `mapstructure:"GORM_LOG_LEVEL" validate:"required,oneof=error warn info silent"`
	POSTGRES_DB_HOST       string `mapstructure:"POSTGRES_DB_HOST" validate:"required"`
	POSTGRES_DB_PORT       string `mapstructure:"POSTGRES_DB_PORT" validate:"required"`
	POSTGRES_DB_USER       string `mapstructure:"POSTGRES_DB_USER" validate:"required"`
	POSTGRES_DB_PASSWORD   string `mapstructure:"POSTGRES_DB_PASSWORD" validate:"required"`
	POSTGRES_DB_NAME       string `mapstructure:"INTERMEDIATOR_POSTGRES_DB_NAME" validate:"required"`
	POSTGRES_DB_SSLMODE    string `mapstructure:"POSTGRES_DB_SSLMODE" validate:"required"`
	JWT_KEY                string `mapstructure:"JWT_KEY" validate:"required"`
	BASIC_AUTH_USERNAME    string `mapstructure:"BASIC_AUTH_USERNAME" validate:"required"`
	BASIC_AUTH_PASSWORD    string `mapstructure:"BASIC_AUTH_PASSWORD" validate:"required"`
	API_KEY                string `mapstructure:"API_KEY" validate:"required"`
	X_API_KEY              string `mapstructure:"X_API_KEY" validate:"required"`
	GESTRYM_SERVER_ADDRESS string `mapstructure:"GESTRYM_SERVER_ADDRESS" validate:"required"`
	EMAIL_CONFIRMATION_URL string `mapstructure:"EMAIL_CONFIRMATION_URL"`
	SUPPORT_EMAIL          string `mapstructure:"SUPPORT_EMAIL"`
}

func (v *env) Validate() error {
	validator := utils.GetValidator()
	structErrors := validator.New().Struct(v)
	if structErrors != nil {
		return structErrors
	}

	// Reflection to iterate over struct fields
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := val.Type().Field(i).Tag.Get("mapstructure")
		if strings.Contains(field.String(), tag) {
			return errors.New("required env variable: " + tag)
		}
	}

	return nil
}

var (
	Env     *env
	envOnce sync.Once
)

const (
	envLocalConfigFile      = ".././deployment/env_local.yaml"
	envProductionConfigFile = "./deployment/env.yaml"
	envTestConfigFile       = "../../../../../deployment/env_test.yaml"
)

func InitEnvironment(isLocalEnv bool) {
	envOnce.Do(func() {
		// Seleccionar el archivo de configuración dependiendo del entorno
		if isLocalEnv {
			logger.Info("[GESTRYM_STORAGE] servidor iniciado en modo local")
			viper.SetConfigFile(envLocalConfigFile)
		} else {
			logger.Info("[GESTRYM_STORAGE] servidor iniciado en modo producción")
		}

		// Leer las variables de entorno
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		bindAllEnvVars()

		// Leer el archivo de configuración
		if isLocalEnv {
			err := viper.ReadInConfig()
			if err != nil {
				logger.Fatal("error reading env file: %v", err)
			}
		}

		// Establecer valores por defecto para las variables de entorno
		setDefaults()

		// Mapear las variables de entorno a la instancia Env
		err := viper.Unmarshal(&Env)
		if err != nil {
			logger.Fatal("unable to decode into struct, %v", err)
		}

		// Validar los datos cargados en la instancia Env
		err = Env.Validate()
		if err != nil {
			logger.Fatal("error validating config: %v", err)
		}

		logger.Success("[OK] variables de entorno cargadas correctamente")
	})
}

func bindAllEnvVars() {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		_ = viper.BindEnv(key)
	}
}

func setDefaults() {
	viper.SetDefault("GIN_MODE", "release")
	viper.SetDefault("GORM_LOG_LEVEL", "error")
}

func InitTestEnvironment() {
	// projectPath, err := utils.GetProjectPath()
	// if err != nil {
	// 	logger.Fatal("error getting project projectPath: %w", err)
	// }
	// logger.Info("[INIT_TEST_ENVIRONMENT_VARIABLES] projectPath: %v", projectPath)

	//envFilePath := projectPath + envTestConfigFile

	logger.Info("[INIT_TEST_ENVIRONMENT_VARIABLES] envFilePath: %v", envTestConfigFile)

	viper.SetConfigFile(envTestConfigFile)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("error reading env file: %v", err)
	}

	err = viper.Unmarshal(&Env)
	if err != nil {
		logger.Fatal("unable to decode into struct, %v", err)
	}

	logger.Success("[INIT_TEST_ENVIRONMENT_VARIABLES] variables de entorno cargadas correctamente")
}
