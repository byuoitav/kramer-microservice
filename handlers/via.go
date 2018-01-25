package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/kramer-microservice/via"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func ResetVia(context echo.Context) error {
	defer color.Unset()
	address := context.Param("address")

	err := via.Reset(address)
	if err != nil {
		color.Set(color.FgRed)
		log.Printf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgGreen, color.Bold)
	log.Printf("Success.")

	return context.JSON(http.StatusOK, "Success")
}

func RebootVia(context echo.Context) error {
	defer color.Unset()
	address := context.Param("address")

	err := via.Reboot(address)
	if err != nil {
		color.Set(color.FgRed)
		log.Printf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgGreen, color.Bold)
	log.Printf("Success.")

	return context.JSON(http.StatusOK, "Success")
}

func SetViaVolume(context echo.Context) error {
	defer color.Unset()
	address := context.Param("address")

	err := via.SetVolume(address)
	if err != nil {
		color.Set(color.FgRed)
		log.Printf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgGreen, color.Bold)
	log.Printf("Success.")

	return context.JSON(http.StatusOK, "Success")
}

func GetViaConnectedStatus(context echo.Context) error {
	address := context.Param("address")

	connected := via.IsConnected(address)

	if connected {
		color.Set(color.FgGreen, color.Bold)
		log.Printf("%s is connected", address)
	} else {
		color.Set(color.FgRed)
		log.Printf("%s is not connected", address)
	}

	return context.JSON(http.StatusOK, connected)
}

func GetViaVolume(context echo.Context) error {
	address := context.Param("address")

	connected := via.GetVolume(address)

	if connected {
		color.Set(color.FgGreen, color.Bold)
		log.Printf("%s is connected", address)
	} else {
		color.Set(color.FgRed)
		log.Printf("%s is not connected", address)
	}

	return context.JSON(http.StatusOK, connected)
}
