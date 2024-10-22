package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mateuszromek/auth/internal/db"
	"github.com/mateuszromek/auth/internal/storage"
	"go.uber.org/zap"
)

type config struct {
	addr   string
	port   string
	dbAddr string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	cfg := config{
		addr:   os.Getenv("APP_ADDRESS"),
		port:   os.Getenv("APP_PORT"),
		dbAddr: os.Getenv("DATABASE_URL"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	logger.Info(cfg.addr, cfg.port)
	db, err := db.NewDb(cfg.dbAddr)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	logger.Info("Database connection established")

	store := storage.NewStorage(db)

	app := &application{
		cfg:    cfg,
		store:  store,
		logger: logger,
	}

	r := app.NewRouter()

	addr := fmt.Sprintf("%s:%s", cfg.addr, cfg.port)
	fmt.Println(addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err)
	}

	logger.Info("Server started on ", addr)
}
