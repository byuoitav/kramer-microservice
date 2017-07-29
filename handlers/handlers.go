package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/kramer-microservice/helpers"
	"github.com/labstack/echo"
)

func SwitchInput(context echo.Context) error {
	log.Printf("Changing inputs")

	input := context.Param("input")
	output := context.Param("output")
	address := context.Param("address")

	log.Printf("Routing %v to %v on %v", input, output, address)

	err := helpers.SwitchInput(address, input, output)
	if err != nil {
		log.Printf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.Printf("Success")
	return context.JSON(http.StatusOK, "Success")
}

func GetCurrentInput(context echo.Context) error {

	//	address := context.Param("address")
	//	inputs, err := helpers.GetCurrentInputs(address)
	//	if err != nil {
	//		return context.JSON(http.StatusInternalServerError, err.Error())
	//	}

	//	return context.JSON(http.StatusOK, inputs)
	return nil
}

func GetInputByPort(context echo.Context) error {

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
	return nil
}

func SetBlank(context echo.Context) error {
	return nil
}
