package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mateuszromek/auth/internal/db"
	"go.uber.org/zap"
)

type config struct {
	addr   string
	dbAddr string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	cfg := config{
		addr:   os.Getenv("APP_ENV"),
		dbAddr: os.Getenv("DATABASE_URL"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.NewDb(cfg.dbAddr)

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	logger.Info("Database connection established")

	fmt.Println(cfg)
}
