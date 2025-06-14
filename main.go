package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/yaml.v3"

	"github.com/burhangltekin/byfood/models"
	"github.com/burhangltekin/byfood/routes"
	"github.com/burhangltekin/byfood/utils"
)

// @title           ByFood API
// @version         1.0
// @description     API for managing books in ByFood.
// @termsOfService  TBD

// @contact.name   API Support
// @contact.url    TBD
// @contact.email  burhangltekin2@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := setupApp(config); err != nil {
		log.Fatalf("Failed to set up app: %v", err)
	}

	r := gin.Default()
	r.Use(gin.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r)

	r.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.URL.Path == "/swagger/doc.json" {
			c.File("./swagger/doc.json")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupApp(config models.AppConfig) error {
	if config.AutoMigrate {
		if err := utils.InitDB(); err != nil {
			log.Fatalf("Database initialization failed: %v", err)
		}
	}
	return nil
}

func loadConfig(path string) (models.AppConfig, error) {
	var config models.AppConfig
	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("Warning: failed to close config file: %v", cerr)
		}
	}()
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}
	return config, nil
}
