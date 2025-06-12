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
		log.Fatal(err)
	}

	err = setupApp(config)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(gin.Logger())

	routes.SetupRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}

func setupApp(config models.AppConfig) error {
	if config.AutoMigrate {
		utils.InitDB()
	}
	return nil
}

func loadConfig(path string) (models.AppConfig, error) {
	var config models.AppConfig
	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	return config, err
}
