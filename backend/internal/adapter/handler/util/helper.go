package utilhandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/handler/response"
	"github.com/gin-gonic/gin"
)

// stringToUint64 is a helper function to convert a string to uint64
func StringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

// getAuthPayload is a helper function to get the auth payload from the context
// func getAuthPayload(ctx *gin.Context, key string) *domain.TokenPayload {
// 	return ctx.MustGet(key).(*domain.TokenPayload)
// }

// toMap is a helper function to add meta and data to a map
func ToMap(m response.Meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}

func SetAuthCookies(ctx *gin.Context, accessToken string, duration time.Duration) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   int(duration),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(ctx.Writer, cookie)
}

func ClearAuthCookies(ctx *gin.Context) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(ctx.Writer, cookie)
}
