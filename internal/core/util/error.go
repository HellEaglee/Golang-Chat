package util

import (
	"errors"
)

var (
	ErrInternal                   = errors.New("internal error")
	ErrDataNotFound               = errors.New("data not found")
	ErrNoUpdatedData              = errors.New("no data to update")
	ErrConflictingData            = errors.New("data conflicts with existing data in unique column")
	ErrInsufficientStock          = errors.New("product stock is not enough")
	ErrInsufficientPayment        = errors.New("total paid is less than total price")
	ErrAccessTokenDuration        = errors.New("invalid access token duration format")
	ErrAccessTokenCreation        = errors.New("error creating access token")
	ErrRefreshTokenDuration       = errors.New("invalid refresh token duration format")
	ErrRefreshTokenCreation       = errors.New("error creating refresh token")
	ErrExpiredAccessToken         = errors.New("access token has expired")
	ErrInvalidAccessToken         = errors.New("access token is invalid")
	ErrExpiredRefreshToken        = errors.New("refresh token has expired")
	ErrInvalidRefreshToken        = errors.New("refresh token is invalid")
	ErrInvalidSession             = errors.New("invalid session id")
	ErrInvalidCredentials         = errors.New("invalid email or password")
	ErrEmptyAuthorizationHeader   = errors.New("authorization header is not provided")
	ErrInvalidAuthorizationHeader = errors.New("authorization header format is invalid")
	ErrInvalidAuthorizationType   = errors.New("authorization type is not supported")
	ErrUnauthorized               = errors.New("user is unauthorized to access the resource")
	ErrForbidden                  = errors.New("user is forbidden to access the resource")
)
