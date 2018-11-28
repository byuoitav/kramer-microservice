package handlers

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/kramer-microservice/via"
	"github.com/labstack/echo"
)

// GetViaHardwareInfo gets the hardware information for a VIA
func GetViaHardwareInfo(context echo.Context) error {
	hardware, err := via.GetHardwareInfo(context.Param("address"))
	if err != nil {
		log.L.Errorf("failed to get the VIA hardware information", err.String())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, hardware)
}

// GetSwitcherHardwareInfo gets the hardware information for a Kramer video switcher
func GetSwitcherHardwareInfo(context echo.Context) error {
	return context.JSON(http.StatusOK, nil)
}
