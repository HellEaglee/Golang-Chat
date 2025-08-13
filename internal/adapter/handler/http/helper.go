package httphandler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// stringToUint64 is a helper function to convert a string to uint64
func stringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

// getAuthPayload is a helper function to get the auth payload from the context
// func getAuthPayload(ctx *gin.Context, key string) *domain.TokenPayload {
// 	return ctx.MustGet(key).(*domain.TokenPayload)
// }

// toMap is a helper function to add meta and data to a map
func toMap(m meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}

func setAuthCookies(ctx *gin.Context, accessToken, refreshToken string) {
	ctx.SetCookie(
		"access_token",
		accessToken,
		15*60, // 15 minutes
		"/",
		"",
		false, // Set to true in production
		true,  // httpOnly
	)

	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60, // 7 days
		"/",
		"",
		false, // Set to true in production
		true,  // httpOnly
	)
}

func clearAuthCookies(ctx *gin.Context) {
	ctx.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}
