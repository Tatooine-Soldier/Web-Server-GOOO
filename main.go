package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HelloWorld struct {
	Message string `json:"message"`
}

type Person struct {
	Name string `json:"name"`
}

func main() {
	e := echo.New()
	e.GET("/", Home)
	e.GET("/contact", Contact)

	e.GET("/params/:data", getParams)
	e.Logger.Fatal(e.Start(":1323"))
}

func getParams(c echo.Context) error {
	person := Person{}

	err := c.Bind(&person)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed processing request")
	}

	datatype := c.Param("data")

	// if datatype == "string" {
	// 	return c.String(http.StatusOK, fmt.Sprintf("This is the param name you sent us '%s'", person.Name))
	// }

	if datatype == "json" {
		return c.String(http.StatusOK, fmt.Sprintf("This is the param name you sent us '%s'", person.Name))
	}

	return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid parameter type: %v", datatype))

}

func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome Home!")
}

func Contact(c echo.Context) error {
	return c.String(http.StatusOK, "Contact!")
}
