package httphandler

import (
	"log/slog"
	"strings"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
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

	v1 := router.Group("/v1")
	{
		post := v1.Group("/posts")
		{
			post.POST("/", postHandler.CreatePost)
		}
	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
