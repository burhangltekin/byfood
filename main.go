package main

import (
	"github.com/burhangltekin/byfood/models"
	"github.com/burhangltekin/byfood/routes"
	"github.com/burhangltekin/byfood/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

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

	routes.SetupRoutes(r)

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
