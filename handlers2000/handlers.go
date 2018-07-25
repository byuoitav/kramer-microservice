package handlers2000

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/byuoitav/common/status"
	"github.com/byuoitav/kramer-microservice/p2000"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func GetInputByPort(context echo.Context) error {
	port := context.Param("port")
	address := context.Param("address")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Printf(color.HiRedString("Need and integer for input and output"))
		return context.JSON(http.StatusBadRequest, "Need and integer for input and output")
	}

	val, err := p2000.GetInputByPort(address, portInt)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "[Kramer 2000 protocol]"+err.Error())
	}

	return context.JSON(http.StatusOK, status.Input{Input: fmt.Sprintf("%v:%v", val, port)})
}

func SwitchInput(context echo.Context) error {
	input := context.Param("input")
	output := context.Param("output")
	address := context.Param("address")

	inputInt, err := strconv.Atoi(input)
	if err != nil {
		log.Printf(color.HiRedString("Need and integer for input and output"))
		return context.JSON(http.StatusBadRequest, "Need and integer for input and output")
	}
	outputInt, err := strconv.Atoi(output)
	if err != nil {
		log.Printf(color.HiRedString("Need and integer for input and output"))
		return context.JSON(http.StatusBadRequest, "Need and integer for input and output")
	}

	val, err := p2000.SetOutput(address, inputInt, outputInt)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "[Kramer 2000 protocol]"+err.Error())
	}

	return context.JSON(http.StatusOK, status.Input{Input: fmt.Sprintf("%v:%v", val, output)})
}
