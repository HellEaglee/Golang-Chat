package httphandler

import (
	"fmt"
	"net/http"

	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
	"github.com/gin-gonic/gin"
)

func authMiddleWare(s port.TokenService, csrf port.CSRFService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !verifyCSRFToken(ctx, csrf) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "CSRF token validation failed"})
			ctx.Abort()
			return
		}

		accessToken, err := ctx.Cookie("access_token")
		if err != nil {
			handleAbort(ctx, err)
			return
		}

		payload, err := s.VerifyToken(accessToken)
		if err != nil {
			if err == util.ErrExpiredAccessToken {
				refreshToken, err := ctx.Cookie("refresh_token")
				if err != nil {
					handleAbort(ctx, err)
					return
				}
				_, err = s.VerifyRefreshToken(ctx, refreshToken)
				if err != nil {
					handleAbort(ctx, err)
					return
				}

				newAccessToken, newRefreshToken, err := s.RefreshTokens(ctx, accessToken, refreshToken)
				if err != nil {
					handleAbort(ctx, err)
					return
				}
				fmt.Printf("refreshed token")
				setAuthCookies(ctx, newAccessToken, newRefreshToken)

				payload, err = s.VerifyToken(newAccessToken)
				if err != nil {
					handleAbort(ctx, err)
					return
				}
			} else {
				handleAbort(ctx, err)
				return
			}
		}

		ctx.Set("user_id", payload.UserID)
		ctx.Next()
	}
}

func verifyCSRFToken(ctx *gin.Context, csrf port.CSRFService) bool {
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
