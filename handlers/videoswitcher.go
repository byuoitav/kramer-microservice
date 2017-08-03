package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	vs "github.com/byuoitav/kramer-microservice/videoswitcher"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func SwitchInput(context echo.Context) error {
	defer color.Unset()

	input := context.Param("input")
	output := context.Param("output")
	address := context.Param("address")

	i, err := vs.ToIndexOne(input)
	if err != nil || vs.LessThanZero(input) {
		return context.JSON(http.StatusBadRequest, fmt.Sprintf("Error! Input parameter %s is not valid!", input))
	}

	o, err := vs.ToIndexOne(output)
	if err != nil || vs.LessThanZero(output) {
		return context.JSON(http.StatusBadRequest, "Error! Output parameter must be zero or greater")
	}

	color.Set(color.FgYellow)
	log.Printf("Routing %v to %v on %v", input, output, address)
	log.Printf("Changing to 1-based indexing... (+1 to each port number)")

	err = vs.SwitchInput(address, i, o)
	if err != nil {
		color.Set(color.FgRed)
		log.Printf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgGreen)
	log.Printf("Success")
	return context.JSON(http.StatusOK, "Success")
}

func GetInputByPort(context echo.Context) error {
	defer color.Unset()

	address := context.Param("address")
	port := context.Param("port")
	p, err := vs.ToIndexOne(port)
	if err != nil || vs.LessThanZero(port) {
		return context.JSON(http.StatusBadRequest, "Error! Port parameter must be zero or greater")
	}

	color.Set(color.FgYellow)
	log.Printf("Getting input for output port %s", port)
	log.Printf("Changing to 1-based indexing... (+1 to each port number)")

	input, err := vs.GetCurrentInputByOutputPort(address, p)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgYellow)
	log.Printf("Changing to 0-based indexing... (-1 to each port number)")
	input.Input, err = vs.ToIndexZero(input.Input)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgGreen)
	log.Printf("Input for output port %s is %v", port, input.Input)
	return context.JSON(http.StatusOK, input)
}

func SetFrontLock(context echo.Context) error {
	defer color.Unset()
	address := context.Param("address")
	state := context.Param("bool")
	stateB, err := strconv.ParseBool(state)
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Error! front-button-lock must be set to true/false")
	}

	color.Set(color.FgYellow)
	log.Printf("Setting front button lock status to %v", stateB)

	err = vs.SetFrontLock(address, stateB)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgGreen)
	log.Printf("Success")
	return context.JSON(http.StatusOK, "Success")
}
