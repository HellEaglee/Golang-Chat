package httphandler

import (
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
