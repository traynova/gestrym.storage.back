package config

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gestrym-storage/src/common/utils"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type IDatabaseConnection interface {
	GetDB() *gorm.DB
	Close() error
}

type postgresConnection struct {
	host       string
	port       string
	user       string
	password   string
	name       string
	sslMode    string
	connection *gorm.DB
	logger     utils.ILogger
}

var (
	dbConnectionInstance *postgresConnection
	databaseOnce         sync.Once
)

func NewPostgresConnection() IDatabaseConnection {
	databaseOnce.Do(func() {
		// Initialize database instance
		dbConnectionInstance = &postgresConnection{}

		// Initialize logger
		dbConnectionInstance.logger = utils.NewLogger()

		// Read environment variables
		err := dbConnectionInstance.readEnvironmentVariables()
		if err != nil {
			dbConnectionInstance.logger.Fatal("failed to read environment variables: %v", err)
		}

		// connect to database
		err = dbConnectionInstance.connect()
		if err != nil {
			dbConnectionInstance.logger.Fatal("failed to connect to database: %v", err)
		}

		// ping database
		err = dbConnectionInstance.ping()
		if err != nil {
			dbConnectionInstance.logger.Fatal("failed to ping database: %v", err)
		}

		dbConnectionInstance.logger.Success("[OK] connected to database")
	})

	// Return database instance
	return dbConnectionInstance
}

func (p *postgresConnection) readEnvironmentVariables() error {
	p.host = viper.GetString("POSTGRES_DB_HOST")
	if p.host == "" {
		return fmt.Errorf("missing required environment variable: POSTGRES_DB_HOST")
	}

	p.port = viper.GetString("POSTGRES_DB_PORT")
	if p.port == "" {
		return fmt.Errorf("missing required environment variable: POSTGRES_DB_PORT")
	}

	p.user = viper.GetString("POSTGRES_DB_USER")
	if p.user == "" {
		return fmt.Errorf("missing required environment variable: POSTGRES_DB_USER")
	}

	p.password = viper.GetString("POSTGRES_DB_PASSWORD")
	if p.password == "" {
		return fmt.Errorf("missing required environment variable: POSTGRES_DB_PASSWORD")
	}

	p.name = viper.GetString("POSTGRES_DB_NAME")
	if p.name == "" {
		return fmt.Errorf("missing required environment variable: POSTGRES_DB_NAME")
	}

	p.sslMode = viper.GetString("POSTGRES_DB_SSLMODE")
	if p.sslMode == "" {
		return fmt.Errorf("missing required environment variable: POSTGRES_DB_SSLMODE")
	}

	return nil
}

func (p *postgresConnection) connect() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", p.host, p.port, p.user, p.password, p.name, p.sslMode)

	var err error
	p.connection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(p.getLoggerLevel()),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}

func (p *postgresConnection) ping() error {
	if p.connection == nil {
		return fmt.Errorf("connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDB, err := p.connection.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (p *postgresConnection) Close() error {
	sqlDB, err := p.connection.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

func (p *postgresConnection) GetDB() *gorm.DB {
	return p.connection
}

func (p *postgresConnection) getLoggerLevel() gormLogger.LogLevel {
	level := strings.ToLower(viper.GetString("GORM_LOG_LEVEL"))
	switch level {
	case "silent":
		return gormLogger.Silent
	case "error":
		return gormLogger.Error
	case "warn":
		return gormLogger.Warn
	case "info":
		return gormLogger.Info
	default:
		return gormLogger.Error
	}
}
