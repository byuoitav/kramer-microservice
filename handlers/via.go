package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	se "github.com/byuoitav/av-api/statusevaluators"
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
	fmt.Printf("Value passed by SetViaVolume is %s", value)

	volume, err := strconv.Atoi(value)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	} else if volume > 100 || volume < 0 {
		return context.JSON(http.StatusBadRequest, "Error: volume must be a value from 0 to 100!")
	}

	volumec := strconv.Itoa(volume)
	log.Printf("Setting volume for %s to %v...", address, volume)

	response, err := via.SetVolume(address, volumec)

	if err != nil {
		log.Printf("An Error Occured: %s", err)
		return context.JSON(http.StatusBadRequest, "An error has occured while setting volume")
	}
	log.Printf("Success: %s", response)
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

	ViaVolume, err := via.GetVolume(address)

	if err != nil {
		color.Set(color.FgRed)
		log.Printf("Failed to retreive VIA volume")
		return context.JSON(http.StatusBadRequest, "Failed to retreive VIA volume")
	} else {
		color.Set(color.FgGreen, color.Bold)
		log.Printf("VIA volume is currently set to %v", strconv.Itoa(ViaVolume))
		return context.JSON(http.StatusOK, se.Volume{ViaVolume})
	}

}
