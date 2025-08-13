package httphandler

import (
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
	"github.com/gin-gonic/gin"
)

func authMiddleWare(s port.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
