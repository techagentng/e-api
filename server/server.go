package server

import (
	"context"
	"fmt"
	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/db"
	"github.com/techagentng/ecommerce-api/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Config         *config.Config
	AuthRepository db.AuthRepository
	AuthService    services.AuthService
	OrderService   services.OrderService
	OrderRepo      db.OrderRepository
	ProductRepo	db.ProductRepository
	DB             db.GormDB
}

// Server serves requests to DB with rout
func (s *Server) Start() {
	r := s.setupRouter()
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if PORT == ":" {
		PORT = ":8080"
	}
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on %s\n", PORT)
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
