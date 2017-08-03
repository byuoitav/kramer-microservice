package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func RestartVIA(context echo.Context) error {
	return context.JSON(http.StatusOK, "Success")
}
