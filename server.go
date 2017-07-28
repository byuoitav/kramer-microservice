package main

import (
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/kramer-microservice/handlers"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := ":8011"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	//Functionality endpoints
	secure.GET("/:address/input/:input/:output", handlers.SwitchInput)
	secure.GET("/:address/buttonlock/:bool", handlers.SetFrontLock)
	secure.GET("/:address/output/:output/blank/:bool", handlers.SetBlank)

	//Status endpoints
	secure.GET("/:address/input/map", handlers.GetCurrentInput)
	secure.GET("/:address/input/get/:port", handlers.GetInputByPort)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
