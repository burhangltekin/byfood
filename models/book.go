package models

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"gte=0,lte=2100"`
}

type BookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"gte=0,lte=2100"`
}

type AppConfig struct {
	LogLevel         string   `yaml:"logLevel"`
	EnableReqLogging bool     `yaml:"enableReqLogging"`
	AutoMigrate      bool     `yaml:"autoMigrate"`
	CORSOrigins      []string `yaml:"corsOrigins"`
	APIVersion       string   `yaml:"apiVersion"`
	ShutdownTimeout  int      `yaml:"shutdownTimeout"`
}
