package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/kramer-microservice/helpers"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func SwitchInput(context echo.Context) error {
	defer color.Unset()

	input := context.Param("input")
	output := context.Param("output")
	address := context.Param("address")

	color.Set(color.FgYellow)
	log.Printf("Routing %v to %v on %v", input, output, address)

	err := helpers.SwitchInput(address, input, output)
	if err != nil {
		color.Set(color.FgRed)
		log.Printf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	color.Set(color.FgGreen)
	log.Printf("Success")
	return context.JSON(http.StatusOK, "Success")
}

func GetCurrentInput(context echo.Context) error {
	defer color.Unset()

	//	address := context.Param("address")
	//	inputs, err := helpers.GetCurrentInputs(address)
	//	if err != nil {
	//		return context.JSON(http.StatusInternalServerError, err.Error())
	//	}

	//	return context.JSON(http.StatusOK, inputs)
	return nil
}

func GetInputByPort(context echo.Context) error {
	defer color.Unset()

	//	address := context.Param("address")
	//	port := context.Param("port")
	//	bay, err := strconv.Atoi(port)
	//	if err != nil || bay < 0 {
	//		return context.JSON(http.StatusBadRequest, "Error! Port parameter must be zero or greater")
	//	}
	//
	//	input, err := helpers.GetInputByOutputPort(address, bay)
	//	if err != nil {
	//		return context.JSON(http.StatusInternalServerError, err.Error())
	//	}
	//
	//	return context.JSON(http.StatusOK, input)
	return nil
}

func SetFrontLock(context echo.Context) error {
	defer color.Unset()
	return nil
}

func SetBlank(context echo.Context) error {
	defer color.Unset()
	return nil
}
