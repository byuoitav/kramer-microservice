package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	var VolumeSetFin string
	address := context.Param("address")
	value := context.Param("volvalue")

	volume, err := strconv.Atoi(value)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	} else if volume > 100 || volume < 0 {
		return context.JSON(http.StatusBadRequest, "Error: volume must be a value from 0 to 100!")
	}
	/*
		if volume > 100 || volume < 0 {
			return context.JSON(http.StatusBadRequest, "Error: volume must be a value from 0 to 100!")
		}
	*/
	log.Printf("Setting volume for %s to %v...", address, volume)

	response, err := via.SetVolume(address, volume)
	/*
		if err != nil {
			color.Set(color.FgRed)
			log.Printf("There was a problem: %v", err.Error())
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	*/
	// Error handling - Handle with care
	// Error1 - value is outside the bounds of 0-100
	// Error2 - no value set for volume
	if strings.Contains(response, "Error1") {
		return context.JSON(http.StatusBadRequest, "Volume command error - volume value is outside the bounds of 0-100")
		log.Printf("Volume command error - volume value %s is outside the bounds of 0-100", "volume")
	} else if strings.Contains(response, "Error2") {
		return context.JSON(http.StatusBadRequest, "volume value was not in the command passed")
		log.Printf("Volume command error - volume value was not in the command passed to %s", address)
	} else {
		//r, _ := regexp.Compile("/\|\d/g")
		r, _ := regexp.Compile("/\\d/g")
		VolumeSetFin = r.FindString(response) //matching strings or should we convert to integer
	}

	if VolumeSetFin != "volume" {
		return context.JSON(http.StatusBadRequest, "Volume command error - volume did not change as requested")
		log.Printf("Volume command error - volume did not change to %s as requested", volume)
	} else {
		log.Printf("Success. Volume changed to %s", volume)
		color.Set(color.FgGreen, color.Bold)

		return context.JSON(http.StatusOK, "Success")
	}
	return context.JSON(http.StatusBadRequest, "An error has occured, try setting volume again")
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

	ViaVolume, _ := via.GetVolume(address)

	if ViaVolume != "" {
		color.Set(color.FgGreen, color.Bold)
		log.Printf("VIA volume is currently set to %s", ViaVolume)
		return context.JSON(http.StatusOK, ViaVolume)
	} else {
		color.Set(color.FgRed)
		log.Printf("Failed to retreive VIA volume")
		return context.JSON(http.StatusBadRequest, "Failed to retreive VIA volume")
	}
}
