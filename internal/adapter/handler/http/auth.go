package httphandler

import (
	"time"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service  port.AuthService
	csrf     port.CSRFService
	duration time.Duration
}

func NewAuthHandler(config *config.Token, service port.AuthService, csrf port.CSRFService) *AuthHandler {
	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		return nil
	}
	return &AuthHandler{service: service, csrf: csrf, duration: duration}
}

type authRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

// Login godoc
//
//	@Summary		Login and get cookies
//	@Description	Logs in a registered user and returns cookies if the credentials are valid.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authRequest		true	"Login request body"
//	@Success		200		{object}	authResponse	"Login successful"
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

	accessToken, err := handler.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}

	setAuthCookies(ctx, accessToken, handler.duration)
	rsp := newAuthResponse("Login successful")

	handleSuccess(ctx, rsp)
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Logs out a registered user and returns access/refresh tokens if the credentials are valid.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	authResponse	"Logout successful"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/auth/logout [get]
func (handler *AuthHandler) Logout(ctx *gin.Context) {
	clearAuthCookies(ctx)
	rsp := newAuthResponse("Logout successful")

	handleSuccess(ctx, rsp)
}

// Register godoc
//
//	@Summary		Register and get an access token
//	@Description	Register a user and returns an access token if the credentials are valid.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authRequest		true	"Register request body"
//	@Success		200		{object}	authResponse	"Register successful"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/auth/register [post]
func (handler *AuthHandler) Register(ctx *gin.Context) {
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

	accessToken, err := handler.service.Register(ctx, user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	setAuthCookies(ctx, accessToken, handler.duration)
	rsp := newAuthResponse("Register successful")

	handleSuccess(ctx, rsp)
}

// CSRF godoc
//
//	@Summary		CSRF
//	@Description	Get csrf token if the credentials are valid.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	csrfResponse	"some token"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/auth/csrf-token [get]
func (handler *AuthHandler) GetCSRFToken(ctx *gin.Context) {
	csrfToken, err := handler.csrf.GenerateToken()
	if err != nil {
		handleError(ctx, util.ErrInternal)
		return
	}

	ctx.SetCookie(
		"csrf_token",
		csrfToken,
		3600,
		"/",
		"",
		false,
		true,
	)
	rsp := newCSRFResponse(csrfToken)
	handleSuccess(ctx, rsp)
}
