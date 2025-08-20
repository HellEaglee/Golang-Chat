package httphandler

import (
	"errors"
	"net/http"
	"time"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var errorStatusMap = map[error]int{
	// Server codes - 5XX
	util.ErrInternal:             http.StatusInternalServerError,
	util.ErrAccessTokenDuration:  http.StatusInternalServerError,
	util.ErrRefreshTokenDuration: http.StatusInternalServerError,
	util.ErrAccessTokenCreation:  http.StatusInternalServerError,
	util.ErrRefreshTokenCreation: http.StatusInternalServerError,

	// Client codes - 4XX
	util.ErrSessionRevoked:  http.StatusGone,
	util.ErrConflictingData: http.StatusConflict,
	util.ErrDataNotFound:    http.StatusNotFound,
	util.ErrNoUpdatedData:   http.StatusBadRequest,

	// Authentication & Authorization code - 401/403
	util.ErrInvalidCredentials:         http.StatusUnauthorized,
	util.ErrUnauthorized:               http.StatusUnauthorized,
	util.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	util.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	util.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	util.ErrInvalidAccessToken:         http.StatusUnauthorized,
	util.ErrInvalidRefreshToken:        http.StatusUnauthorized,
	util.ErrExpiredAccessToken:         http.StatusUnauthorized,
	util.ErrExpiredRefreshToken:        http.StatusUnauthorized,
	util.ErrInvalidSession:             http.StatusUnauthorized,
	util.ErrForbidden:                  http.StatusForbidden,
}

type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

type meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

func newMeta(total, limit, skip uint64) meta {
	return meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

type authResponse struct {
	Message string `json:"message"`
}

func newAuthResponse(message string) authResponse {
	return authResponse{Message: message}
}

type userResponse struct {
	ID        uuid.UUID `json:"id" example:"3342a227-1f2d-4422-a718-435c6a115f62"`
	Name      string    `json:"name" example:"John"`
	Email     string    `json:"email" example:"john@gmail.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

type csrfResponse struct {
	CSRFToken string `json:"csrf_token"`
}

func newCSRFResponse(message string) csrfResponse {
	return csrfResponse{CSRFToken: message}
}

func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
