package httphandler

import (
	"strconv"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(service port.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		createUserRequest	true	"Create user request"
//	@Success		201		{object}	userResponse		"User created"
//	@Failure		400		{object}	errorResponse		"Validation error"
//	@Failure		401		{object}	errorResponse		"Unauthorized error"
//	@Failure		404		{object}	errorResponse		"Data not found error"
//	@Failure		409		{object}	errorResponse		"Data conflict error"
//	@Failure		500		{object}	errorResponse		"Internal server error"
//	@Router			/users [post]
func (handler *UserHandler) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := &domain.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := handler.service.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(createdUser)
	handleSuccess(ctx, rsp)
}

type getUserRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetUser godoc
//
//	@Summary		Get a user by ID
//	@Description	Get a single user by its UUID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"User ID (UUID)"
//	@Success		200	{object}	userResponse	"User found"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users/{id} [get]
func (handler *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, err := handler.service.GetUser(ctx.Request.Context(), req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(user)

	handleSuccess(ctx, rsp)
}

type getUsersRequest struct {
	Skip  string `form:"skip" binding:"required,numeric" example:"0"`
	Limit string `form:"limit" binding:"required,numeric,min=1" example:"5"`
}

// GetUsers godoc
//
//	@Summary		List users with pagination
//	@Description	Get a paginated list of users
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		int				true	"Number of items to skip"	example(0)
//	@Param			limit	query		int				true	"Number of items to take"	example(5)	minimum(1)
//	@Success		200		{object}	meta			"Users displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/users [get]
func (handler *UserHandler) GetUsers(ctx *gin.Context) {
	var req getUsersRequest
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

	users, err := handler.service.GetUsers(ctx.Request.Context(), skip, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	userResponses := make([]userResponse, len(users))
	for i, user := range users {
		userResponses[i] = newUserResponse(&user)
	}

	total := uint64(len(users))
	meta := newMeta(total, limit, skip)
	rsp := toMap(meta, userResponses, "users")

	handleSuccess(ctx, rsp)
}

type updateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,required"`
	Email    string `json:"email" binding:"omitempty,required"`
	Password string `json:"password" binding:"omitempty,min=8"`
}

// UpdateUser godoc
//
//	@Summary		Update an user
//	@Description	Update an existing user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"User ID (UUID)"
//	@Param			user	body		updateUserRequest	true	"Fields to update"
//	@Success		200		{object}	userResponse		"User updated"
//	@Failure		400		{object}	errorResponse		"Validation error"
//	@Failure		401		{object}	errorResponse		"Unauthorized error"
//	@Failure		403		{object}	errorResponse		"Forbidden error"
//	@Failure		404		{object}	errorResponse		"Data not found error"
//	@Failure		500		{object}	errorResponse		"Internal server error"
//	@Router			/users/{id} [put]
func (handler *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
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

	user := &domain.User{
		ID:    uuid,
		Name:  req.Name,
		Email: req.Email,
	}

	if req.Password != "" {
		user.Password = req.Password
	}

	updatedUser, err := handler.service.UpdateUser(ctx.Request.Context(), user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newUserResponse(updatedUser)

	handleSuccess(ctx, rsp)
}

type deleteUserRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// DeletePost godoc
//
//	@Summary		Delete a user
//	@Description	Delete an user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Post ID (UUID)"
//	@Success		200	{object}	response		"User deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users/{id} [delete]
func (handler *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := handler.service.DeleteUser(ctx.Request.Context(), req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
