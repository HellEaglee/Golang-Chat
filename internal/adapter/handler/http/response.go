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

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	util.ErrInternal:                   http.StatusInternalServerError,
	util.ErrDataNotFound:               http.StatusNotFound,
	util.ErrConflictingData:            http.StatusConflict,
	util.ErrInvalidCredentials:         http.StatusUnauthorized,
	util.ErrUnauthorized:               http.StatusUnauthorized,
	util.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	util.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	util.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	util.ErrInvalidAccessToken:         http.StatusUnauthorized,
	util.ErrExpiredAccessToken:         http.StatusUnauthorized,
	util.ErrForbidden:                  http.StatusForbidden,
	util.ErrNoUpdatedData:              http.StatusBadRequest,
	util.ErrInsufficientStock:          http.StatusBadRequest,
	util.ErrInsufficientPayment:        http.StatusBadRequest,
}

type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// meta represents metadata for a paginated response
type meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// newMeta is a helper function to create metadata for a paginated response
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

// newPostResponse is a helper function to create a response body for handling post data
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// handleAbort sends an error response and aborts the request with the specified status code and error message
func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
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

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
