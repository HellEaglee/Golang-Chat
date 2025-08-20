package httphandler

import (
	"log/slog"
	"strings"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter(config *config.HTTP, tokenConfig *config.Token,
	token port.TokenService, csrf port.CSRFService, authHandler AuthHandler, userHandler UserHandler,
) (*Router, error) {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originList
	ginConfig.AllowCredentials = true

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/logout", authHandler.Logout)
			auth.POST("/register", authHandler.Register)
			auth.GET("/csrf-token", authHandler.GetCSRFToken)
		}
		users := v1.Group("/users")
		users.Use(authMiddleWare(token, csrf, tokenConfig))
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/profile", userHandler.GetProfile)
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
