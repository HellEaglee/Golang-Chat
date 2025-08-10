package httphandler

import (
	"strconv"

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

type createPostRequest struct {
	Title       string `json:"title" binding:"required" example:"PostTitleExample"`
	Description string `json:"description" binding:"required" example:"PostDescriptionExample"`
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new blog post with title and description
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		createPostRequest	true	"Create post request"
//	@Success		201		{object}	postResponse		"Post created"
//	@Failure		400		{object}	errorResponse		"Validation error"
//	@Failure		401		{object}	errorResponse		"Unauthorized error"
//	@Failure		404		{object}	errorResponse		"Data not found error"
//	@Failure		409		{object}	errorResponse		"Data conflict error"
//	@Failure		500		{object}	errorResponse		"Internal server error"
//	@Router			/posts [post]
func (handler *PostHandler) CreatePost(ctx *gin.Context) {
	var req createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	post := &domain.Post{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
	}

	createdPost, err := handler.service.CreatePost(ctx.Request.Context(), post)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newPostResponse(createdPost)
	handleSuccess(ctx, rsp)
}

type getPostRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetPost godoc
//
//	@Summary		Get a post by ID
//	@Description	Get a single post by its UUID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Post ID (UUID)"
//	@Success		200	{object}	postResponse	"Post found"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/posts/{id} [get]
func (handler *PostHandler) GetPost(ctx *gin.Context) {
	var req getPostRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	post, err := handler.service.GetPost(ctx.Request.Context(), req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newPostResponse(post)

	handleSuccess(ctx, rsp)
}

type getPostsRequest struct {
	Skip  string `form:"skip" binding:"required,numeric" example:"0"`
	Limit string `form:"limit" binding:"required,numeric,min=1" example:"5"`
}

// GetPosts godoc
//
//	@Summary		List posts with pagination
//	@Description	Get a paginated list of posts
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		int				true	"Number of items to skip"	example(0)
//	@Param			limit	query		int				true	"Number of items to take"	example(5)	minimum(1)
//	@Success		200		{object}	meta			"Posts displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/posts [get]
func (handler *PostHandler) GetPosts(ctx *gin.Context) {
	var req getPostsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	skip, err := strconv.ParseUint(req.Skip, 10, 64)
	if err != nil {
		handleError(ctx, err)
		return
	}

	limit, err := strconv.ParseUint(req.Limit, 10, 64)
	if err != nil {
		handleError(ctx, err)
		return
	}

	posts, err := handler.service.GetPosts(ctx.Request.Context(), skip, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	total := uint64(len(posts))
	meta := newMeta(total, limit, skip)
	rsp := toMap(meta, posts, "posts")

	handleSuccess(ctx, rsp)
}

type updatePostRequest struct {
	Title       string `json:"title" binding:"omitempty,required"`
	Description string `json:"description" binding:"omitempty,required"`
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update an existing post by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Post ID (UUID)"
//	@Param			post	body		updatePostRequest	true	"Fields to update"
//	@Success		200		{object}	postResponse		"Post updated"
//	@Failure		400		{object}	errorResponse		"Validation error"
//	@Failure		401		{object}	errorResponse		"Unauthorized error"
//	@Failure		403		{object}	errorResponse		"Forbidden error"
//	@Failure		404		{object}	errorResponse		"Data not found error"
//	@Failure		500		{object}	errorResponse		"Internal server error"
//	@Router			/posts/{id} [put]
func (handler *PostHandler) UpdatePost(ctx *gin.Context) {
	var req updatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		validationError(ctx, err)
		return
	}

	post := &domain.Post{
		ID:          uuid,
		Title:       req.Title,
		Description: req.Description,
	}

	updatedPost, err := handler.service.UpdatePost(ctx.Request.Context(), post)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newPostResponse(updatedPost)

	handleSuccess(ctx, rsp)
}

type deletePostRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// DeletePost godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by id
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Post ID (UUID)"
//	@Success		200	{object}	response		"User deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/posts/{id} [delete]
func (handler *PostHandler) DeletePost(ctx *gin.Context) {
	var req deletePostRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := handler.service.DeletePost(ctx.Request.Context(), req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
