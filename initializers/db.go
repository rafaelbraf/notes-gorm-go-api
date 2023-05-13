package initializers

import (
	"fmt"
	"log"
	"os"	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"notes/models"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar ao Banco de Dados! \n", err.Error())
		os.Exit(1)
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Rodando migrações...")
	DB.AutoMigrate(&models.Note{})

	log.Println("Conectado ao Banco de Dados com sucesso!")
}