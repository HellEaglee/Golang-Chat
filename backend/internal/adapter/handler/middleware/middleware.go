package middleware

import (
	"errors"
	"time"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/handler/response"
	utilhandler "github.com/HellEaglee/Golang-Chat/internal/adapter/handler/util"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(s port.TokenService, csrf port.CSRFService, config *config.Token) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		duration, err := time.ParseDuration(config.Duration)
		if err != nil {
			response.HandleAbort(ctx, err)
			return
		}
		// if !verifyCSRFToken(ctx, csrf) {
		// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "CSRF token validation failed"})
		// 	ctx.Abort()
		// 	return
		// }

		accessToken, err := ctx.Cookie("access_token")
		if err != nil {
			response.HandleAbort(ctx, util.ErrInvalidAccessToken)
			return
		}

		payload, err := s.VerifyToken(accessToken)
		if err != nil {
			if errors.Is(err, util.ErrExpiredAccessToken) {
				claims, err := s.ExtractClaimsFromToken(accessToken)
				if err != nil {
					response.HandleAbort(ctx, err)
					return
				}
				refreshToken, err := s.GetTokenBySessionID(ctx, claims.SessionID)
				if err != nil {
					response.HandleAbort(ctx, err)
					return
				}
				_, err = s.VerifyRefreshToken(ctx, refreshToken.Token)
				if err != nil {
					response.HandleAbort(ctx, err)
					return
				}

				newAccessToken, err := s.RefreshTokens(ctx, accessToken, refreshToken.Token)
				if err != nil {
					response.HandleAbort(ctx, err)
					return
				}
				utilhandler.SetAuthCookies(ctx, newAccessToken, duration)

				payload, err = s.VerifyToken(newAccessToken)
				if err != nil {
					response.HandleAbort(ctx, err)
					return
				}
			} else {
				response.HandleAbort(ctx, err)
				return
			}
		}

		ctx.Set("user_id", payload.UserID)
		ctx.Next()
	}
}

func VerifyCSRFToken(ctx *gin.Context, csrf port.CSRFService) bool {
	csrfToken := ctx.GetHeader("X-CSRF-TOKEN")
	if csrfToken == "" {
		csrfToken = ctx.PostForm("csrf_token")
	}
	expectedToken, err := ctx.Cookie("csrf_token")
	if err != nil {
		return false
	}
	return csrf.VerifyToken(csrfToken, expectedToken)
}
