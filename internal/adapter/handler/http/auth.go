package httphandler

import (
	"fmt"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service port.AuthService
}

func NewAuthHandler(service port.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type authRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

// Login godoc
//
//	@Summary		Login and get access/refresh tokens
//	@Description	Logs in a registered user and returns access/refresh tokens if the credentials are valid.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authRequest	true	"Login request body"
//	@Success		200		{object}	authResponse	"Succesfully logged in"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/auth/login [post]
func (handler *AuthHandler) Login(ctx *gin.Context) {
	var req authRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	accessToken, refreshToken, err := handler.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAuthResponse(accessToken, refreshToken)

	handleSuccess(ctx, rsp)
}

// Register godoc
//
//	@Summary		Register and get an access token
//	@Description	Register a user and returns an access token if the credentials are valid.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authRequest	true	"Register request body"
//	@Success		200		{object}	authResponse	"Succesfully logged in"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/auth/register [post]
func (handler *AuthHandler) Register(ctx *gin.Context) {
	fmt.Println("Register endpoint hit") // Add this for debugging
	var req authRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := &domain.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
	}

	accessToken, refreshToken, err := handler.service.Register(ctx, user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAuthResponse(accessToken, refreshToken)

	handleSuccess(ctx, rsp)
}
