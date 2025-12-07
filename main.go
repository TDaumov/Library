package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"library_vebservice/handler"
	"library_vebservice/logger"
	"library_vebservice/models"
	"library_vebservice/repository"
	"library_vebservice/service"
	"os"
)

func createAdminIfNotExists(db *gorm.DB, log zerolog.Logger) {
	adminEmail := "admin@example.com"
	adminPassword := "admin123" // поменяй на свой безопасный пароль

	var user models.User
	if err := db.Where("email = ?", adminEmail).First(&user).Error; err == gorm.ErrRecordNotFound {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		admin := models.User{
			Email:    adminEmail,
			Password: string(hashedPassword),
			Role:     "admin",
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create admin user")
		} else {
			log.Info().Str("email", adminEmail).Msg("Admin user created successfully")
		}
	} else if err != nil {
		log.Error().Err(err).Msg("Failed to query admin user")
	} else {
		log.Info().Str("email", adminEmail).Msg("Admin user already exists")
	}
}

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DB_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")

	log := logger.InitLogger() // инициализация логгера

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect db")
	}

	if err := db.AutoMigrate(&models.User{}, &models.Book{}, &models.Borrow{}); err != nil {
		log.Fatal().Err(err).Msg("migration failed")
	}

	if err := db.AutoMigrate(&models.User{}, &models.Book{}, &models.Borrow{}); err != nil {
		log.Fatal().Err(err).Msg("migration failed")
	}

	// Repositories
	userRepo := repository.NewUserRepo(db)
	bookRepo := repository.NewBookRepo(db)
	// Создание админа

	createAdminIfNotExists(db, log)

	// Services
	authSvc := service.NewAuthService(userRepo)
	bookSvc := service.NewBookService(bookRepo)

	//log := logger.InitLogger() // инициализация логгера

	log.Info().Msg("Starting application...")
	log.Error().Msg("This is an error message")
	// Handler
	h := handler.NewHandler(authSvc, bookSvc, log, jwtSecret)

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Регистрация маршрутов
	h.RegisterRoutes(e)

	addr := fmt.Sprintf(":%s", port)
	log.Info().Msgf("Starting server on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
