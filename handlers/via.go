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
	value := context.Param("volvalue")

	ViaVolumeLevel, err := strconv.Atoi(value)
	if err != nil {
			return context.JSON(http.StatusBadRequest, err.Error())
		} else if volume > 100 || volume < 0 {
			return context.JSON(http.StatusBadRequest, "Error: volume must be a value from 0 to 100!")
		}

	log.Printf("Setting volume for %s to %v...", address, volume)

	response, err := via.SetVolume(address,volume)
/*
	if err != nil {
		color.Set(color.FgRed)
		log.Printf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
*/
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

	ViaVolume := via.GetVolume(address,volcurrentlevel)

	if ViaVolume != nil {
		color.Set(color.FgGreen, color.Bold)
		log.Printf("VIA volume is currently set to %s", ViaVolume)
		return context.JSON(http.StatusOK, ViaVolume)
	} else {
		color.Set(color.FgRed)
		log.Printf("Failed to retreive VIA volume")
		return context.JSON(http.StatusBadRequest)
	}

}
