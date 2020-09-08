package main

import (
	"github.com/labstack/echo/v4"

	"app/internal/auth"
	"app/internal/boards"
	"app/internal/items"
	"app/internal/middleware"
	"app/internal/sessionstore"
	"app/internal/validator"
)

func main() {
	e := echo.New()
	e.Validator = validator.NewValidator()

	m := e.Group("", middleware.SessionMiddleware)

	m.GET("/v1/boards", boards.GetBoards())
	m.POST("/v1/boards", boards.AddBoard())

	m.GET("/v1/boards/:id", boards.GetBoard())
	m.PUT("/v1/boards/:id", boards.UpdateBoard())
	m.DELETE("/v1/boards/:id", boards.DeleteBoard())

	m.GET("/v1/boards/:code/id", boards.GetBoardIdByCode())

	m.GET("/v1/boards/:id/members", boards.GetBoardMembers())
	m.POST("/v1/boards/:id/members", boards.JoinBoard())
	m.DELETE("/v1/boards/:id/members/:uid", boards.DeleteBoardMember())

	m.GET("/v1/boards/:id/items", items.GetBoardItems())
	m.POST("/v1/boards/:id/items", items.AddBoardItem())

	m.GET("/v1/boards/:id/items/:iid", items.GetBoardItem())
	m.PUT("/v1/boards/:id/items/:iid", items.UpdateBoardItem())
	m.DELETE("/v1/boards/:id/items/:iid", items.DeleteBoardItem())

	e.POST("/v1/auth/token", auth.GetCsrfToken())
	e.POST("/v1/auth/signin", auth.Signin())
	e.GET("/v1/auth/signout", auth.Signout())

	// e.Startの中はdocker-composeのgoコンテナで設定したportsを指定してください。
	e.Logger.Fatal(e.Start(":8082"))

	defer sessionstore.CloseInstance()
}
