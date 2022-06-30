package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func queryParams(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	return c.String(http.StatusOK, fmt.Sprintf("your cat name is : %s\nand cat type is : %s\n", catName, catType))
}
