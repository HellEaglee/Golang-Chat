package http

import (
	"fmt"
	"strconv"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/handler/response"
	utilhandler "github.com/HellEaglee/Golang-Chat/internal/adapter/handler/util"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
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
//	@Param			user	body		createUserRequest		true	"Create user request"
//	@Success		201		{object}	response.UserResponse	"User created"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		404		{object}	response.ErrorResponse	"Data not found error"
//	@Failure		409		{object}	response.ErrorResponse	"Data conflict error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users [post]
func (handler *UserHandler) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	user := &domain.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := handler.service.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewUserResponse(createdUser)
	response.HandleSuccess(ctx, rsp)
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
//	@Param			id	path		string					true	"User ID (UUID)"
//	@Success		200	{object}	response.UserResponse	"User found"
//	@Failure		400	{object}	response.ErrorResponse	"Validation error"
//	@Failure		404	{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users/{id} [get]
func (handler *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	user, err := handler.service.GetUser(ctx.Request.Context(), req.ID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewUserResponse(user)

	response.HandleSuccess(ctx, rsp)
}

// GetProfile godoc
//
//	@Summary		Get profile data by ID
//	@Description	Get profile data by its UUID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.UserResponse	"Profile found"
//	@Failure		400	{object}	response.ErrorResponse	"Validation error"
//	@Failure		404	{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users/profile [get]
func (handler *UserHandler) GetProfile(ctx *gin.Context) {
	userIDI, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, util.ErrForbidden)
		return
	}

	userID, ok := userIDI.(uuid.UUID)
	if !ok {
		response.HandleError(ctx, fmt.Errorf("invalid user id format"))
		return
	}
	user, err := handler.service.GetUser(ctx.Request.Context(), userID.String())
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewUserResponse(user)

	response.HandleSuccess(ctx, rsp)
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
//	@Param			skip	query		int						true	"Number of items to skip"	example(0)
//	@Param			limit	query		int						true	"Number of items to take"	example(5)	minimum(1)
//	@Success		200		{object}	response.Meta			"Users displayed"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users [get]
func (handler *UserHandler) GetUsers(ctx *gin.Context) {
	var req getUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	skip, err := strconv.ParseUint(req.Skip, 10, 64)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	limit, err := strconv.ParseUint(req.Limit, 10, 64)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	users, err := handler.service.GetUsers(ctx.Request.Context(), skip, limit)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	userResponses := make([]response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = response.NewUserResponse(&user)
	}

	total := uint64(len(users))
	meta := response.NewMeta(total, limit, skip)
	rsp := utilhandler.ToMap(meta, userResponses, "users")

	response.HandleSuccess(ctx, rsp)
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
//	@Param			id		path		string					true	"User ID (UUID)"
//	@Param			user	body		updateUserRequest		true	"Fields to update"
//	@Success		200		{object}	response.UserResponse	"User updated"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404		{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users/{id} [put]
func (handler *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		response.ValidationError(ctx, err)
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
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewUserResponse(updatedUser)

	response.HandleSuccess(ctx, rsp)
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
//	@Param			id	path		string					true	"Post ID (UUID)"
//	@Success		200	{object}	response.Response		"User deleted"
//	@Failure		400	{object}	response.ErrorResponse	"Validation error"
//	@Failure		401	{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403	{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404	{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users/{id} [delete]
func (handler *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := handler.service.DeleteUser(ctx.Request.Context(), req.ID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}
