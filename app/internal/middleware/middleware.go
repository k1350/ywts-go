package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"app/internal/sessionstore"
)

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		aSession, err := sessionstore.GetInstance().Store.Get(c.Request(), "auth-session")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		csrfToken := aSession.Values["csrfToken"]
		if csrfToken == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		uid := aSession.Values["uid"]
		if uid == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		c.Set("uid", uid.(string))
		c.Set("csrfToken", csrfToken.(string))
		return next(c)
	}
}

func CheckSession(c echo.Context) bool {
	csrfToken := c.Get("csrfToken")
	tokenHeader := c.Request().Header.Get("X-CSRF-Token")
	if tokenHeader == "" {
		log.Print("token not found.")
		return false
	}
	if tokenHeader != csrfToken {
		log.Print("session error.")
		return false
	}
	return true
}
