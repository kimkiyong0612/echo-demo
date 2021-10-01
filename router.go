package main

import (
	"echo-demo/api/web"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewRouter(api *web.API) *echo.Echo {

	e := echo.New()

	// Debug mode
	// e.Debug = true

	// Middleware TODO:auth
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(web.WithSessionUser)

	// Health check
	e.GET("/health", healthCheckHandler)

	// Routes
	// TODO: versioning
	e.GET("v1/users", api.GetAllUser)
	e.POST("v1/users", api.CreateUser)
	e.GET("v1/users/:id", api.GetUser)
	e.PUT("v1/users/:id", api.UpdateUser)
	e.DELETE("v1/users/:id", api.DeleteUser)

	return e
}

func healthCheckHandler(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	return ctx.String(http.StatusOK, "{alive:true}")
}
