package httphandler

import (
	"log/slog"
	"strings"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter(config *config.HTTP, postHandler PostHandler) (*Router, error) {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
		posts := v1.Group("/posts")
		{
			posts.POST("/", postHandler.CreatePost)
			posts.GET("/", postHandler.GetPosts)
			posts.GET("/:id", postHandler.GetPost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
		}
	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
