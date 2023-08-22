package app

import (

	// "live_attendance/main/app/database/seeders"
	"live_attendance/main/app/controllers"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Run() {
	var server = controllers.Server{}
	var appConfig = controllers.AppConfig{}
	var dbConfig = controllers.DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error on loading .env file")
	}

	appConfig.AppName = getEnv("APP_NAME", "Live Attendance App")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBUser = getEnv("DB_USER", "root")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "")
	dbConfig.DBName = getEnv("DB_NAME", "live_attendance_db")
	dbConfig.DBPort = getEnv("DB_PORT", "3306")
	dbConfig.DBDriver = getEnv("DB_DRIVER", "mysql")

	server.Initialize(appConfig, dbConfig)
	server.Run(":" + appConfig.AppPort)
}
