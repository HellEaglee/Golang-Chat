package httphandler

import (
	"net/http"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	service port.PostService
}

func NewPostHandler(service port.PostService) *PostHandler {
	return &PostHandler{service: service}
}

type CreatePostRequest struct {
	Title       string `json:"title" binding:"required" example:"PostTitleExample"`
	Description string `json:"description" binding:"required" example:"PostDescriptionExample"`
}

func (handler *PostHandler) CreatePost(ctx *gin.Context) {
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &domain.Post{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
	}

	createPost, err := handler.service.CreatePost(ctx.Request.Context(), post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"post": createPost})
}
