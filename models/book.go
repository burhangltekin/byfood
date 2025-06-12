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
	LogLevel         string
	EnableReqLogging bool
	AutoMigrate      bool
	CORSOrigins      []string
	APIVersion       string
	ShutdownTimeout  int
}
