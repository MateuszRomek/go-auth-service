package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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

	fmt.Println(cfg)
}
