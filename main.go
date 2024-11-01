package main

import (
	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/db"
	"github.com/techagentng/ecommerce-api/server"
	"github.com/techagentng/ecommerce-api/services"
	"log"
	_ "net/url"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	authService := services.NewAuthService(authRepo, conf)

	s := &server.Server{
		Config:         conf,
		AuthRepository: authRepo,
		AuthService:    authService,
		DB:             db.GormDB{},
	}

	s.Start()
}
