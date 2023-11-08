package routes

import echo "github.com/labstack/echo/v4"

func AuthGroup(authG *echo.Group) {
	authG.POST("/login", func(ctx echo.Context) error { return nil })
}
