package main

import (
	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/db"
	"github.com/techagentng/ecommerce-api/server"
	"github.com/techagentng/ecommerce-api/services"
	 "github.com/techagentng/ecommerce-api/docs"
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
	orderRepo := db.NewOrderRepo(gormDB)
	productRepo := db.NewProductRepo(gormDB)
	authService := services.NewAuthService(authRepo, conf)
	orderService := services.NewOrderService(orderRepo, conf)

	s := &server.Server{
		Config:         conf,
		AuthRepository: authRepo,
        OrderRepo: orderRepo,
		AuthService:    authService,
		OrderService:   orderService,
		ProductRepo: productRepo,
		DB:             db.GormDB{},
	}

	s.Start()
}
