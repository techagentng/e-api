package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/swaggo/gin-swagger"
	//  "github.com/swaggo/gin-swagger/swaggerFiles"
	
)

// @title           E-Commerce assessment
// @version         1.0
// @description     Test go skill
// @host            localhost:8080
// @BasePath        /api/v1
func (s *Server) setupRouter() *gin.Engine {
	// ginMode := os.Getenv("GIN_MODE")
	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.MaxMultipartMemory = 32 << 20
	s.defineRoutes(r)

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (s *Server) defineRoutes(router *gin.Engine) {
	apirouter := router.Group("/api/v1")
	apirouter.POST("/auth/signup", s.handleSignup())
	apirouter.POST("/auth/login", s.handleLogin())

	authorized := apirouter.Group("/")
	authorized.Use(s.Authorize())

	// Define user-related routes
	authorized.POST("/user/place/order", s.handlePlaceOrder())
	authorized.GET("/user/orders", s.handleListUserOrders())
	authorized.PATCH("/cancel/order/:order_id", s.handleCancelOrder())
	authorized.PATCH("/update/order/:order_id", s.handleUpdateOrderStatus())
	authorized.POST("/products", s.handleCreateProduct())
	authorized.GET("/products/:product_id", s.handleReadProduct())
	authorized.PUT("/products/:product_id", s.handleUpdateProduct())
	authorized.DELETE("/products/:product_id", s.handleDeleteProduct())
}
