package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/mateuszromek/auth/internal/storage"
	"go.uber.org/zap"
)

type application struct {
	cfg       config
	store     storage.Storage
	logger    *zap.SugaredLogger
	validator *validator.Validate
}
