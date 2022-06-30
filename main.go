package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HelloWorld struct {
	Message string `json:"message"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/contact", func(c echo.Context) error {
		return c.String(http.StatusOK, "Contact!")
	})

	e.GET("/contact/:name", GreetingsWithParams)
	e.Logger.Fatal(e.Start(":1323"))
}

func GreetingsWithParams(c echo.Context) error {
	params := c.Param("name")
	return c.JSON(http.StatusOK, HelloWorld{
		Message: "Hello World, my name is " + params,
	})
}
