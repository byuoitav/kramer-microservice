package videoswitcher

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"

	"github.com/byuoitav/common/status"
)

// Command constants
const (
	BuildDate       = "BUILD-DATE"
	Model           = "MODEL"
	SerialNumber    = "SN"
	FirmwareVersion = "VERSION"
	ProtocolVersion = "PROT-VER"
	Temperature     = "HW-TEMP"
	PowerSave       = "POWER-SAVE"
	IPAddress       = "NET-IP"
	Gateway         = "NET-GATE"
	MACAddress      = "NET-MAC"
	NetDNS          = "NET-DNS"
	Signal          = "SIGNAL"
)

// GetCurrentInputByOutputPort gets the current input that is set to the given output port
func GetCurrentInputByOutputPort(address, port string, readWelcome bool) (status.Input, error) {
	command := fmt.Sprintf("#VID? %s", port)

	respChan := make(chan Response)

	c := CommandInfo{respChan, address, command, readWelcome}

	StartChannel <- c

	re := <-respChan

	resp := re.Response
	err := re.Err

	if err != nil {
		logError(err.Error())
		return status.Input{}, err
	}

	if strings.Contains(resp, "VID") {
		parts := strings.Split(resp, "VID")
		resp = strings.TrimSpace(parts[1])

		parts = strings.Split(resp, ">")

		var i status.Input
		i.Input = parts[0]
		return i, nil
	}

	logError(fmt.Sprintf("Incorrect response for command (%s). (Response: %s)", command, resp))
	return status.Input{}, fmt.Errorf("Incorrect response for command (%s). (Response: %s)", command, resp)
}

// GetHardwareInformation builds the list of hardware information for this device and returns it
func GetHardwareInformation(address string, readWelcome bool) (structs.HardwareInfo, *nerr.E) {
	var toReturn structs.HardwareInfo

	// get the hostname
	addr, e := net.LookupAddr(address)
	if e != nil {
		toReturn.Hostname = address
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
	}

	// get build date
	buildDate, err := hardwareCommand(BuildDate, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get build date from %s", address)
	}

	toReturn.BuildDate = buildDate

	// get device model
	model, err := hardwareCommand(Model, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get model number from %s", address)
	}

	toReturn.ModelName = model

	// get device protocol version
	protocol, err := hardwareCommand(ProtocolVersion, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get protocol version from %s", address)
	}

	toReturn.ProtocolVersion = strings.Trim(protocol, "3000:")

	// get firmware version
	firmware, err := hardwareCommand(FirmwareVersion, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get firmware version from %s", address)
	}

	toReturn.FirmwareVersion = firmware

	// get serial number
	serial, err := hardwareCommand(SerialNumber, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get serial number from %s", address)
	}

	toReturn.SerialNumber = serial

	// get temperature
	// temp, err := hardwareCommand(Temperature, "0", address, readWelcome)
	// if err != nil {
	// 	return toReturn, nerr.Translate(err).Addf("failed to get temperature from %s", address)
	// }

	// toReturn.Temperature, _ = strconv.Atoi(temp)
	// toReturn.Temperature = temp

	// get power saving mode status
	// powerSave, err := hardwareCommand(PowerSave, "", address, readWelcome)
	// if err != nil {
	// 	return toReturn, nerr.Translate(err).Addf("failed to get power saving mode status from %s", address)
	// }

	// toReturn.PowerSavingModeStatus = powerSave

	// get IP address
	ipAddress, err := hardwareCommand(IPAddress, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get IP address from %s... ironic...", address)
	}

	// get gateway
	gateway, err := hardwareCommand(Gateway, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the gateway address from %s", address)
	}

	// get MAC address
	mac, err := hardwareCommand(MACAddress, "", address, readWelcome)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the MAC address from %s", address)
	}

	// get DNS address(es?)
	// dns, err := hardwareCommand(NetDNS, "", address, readWelcome)
	// if err != nil {
	// 	return toReturn, nerr.Translate(err).Addf("failed to get the DNS addresses from %s", address)
	// }

	// set network information
	toReturn.NetworkInfo = structs.NetworkInfo{
		IPAddress:  ipAddress,
		MACAddress: mac,
		Gateway:    gateway,
		// DNS:        []string{dns},
	}

	return toReturn, nil
}

func hardwareCommand(commandType, param, address string, readWelcome bool) (string, error) {
	var command string

	if len(param) > 0 {
		num, _ := strconv.Atoi(param)
		command = fmt.Sprintf("#%s? %d", commandType, num)
	} else {
		command = fmt.Sprintf("#%s?", commandType)
	}

	respChan := make(chan Response)

	c := CommandInfo{respChan, address, command, readWelcome}

	StartChannel <- c

	re := <-respChan

	resp := re.Response
	err := re.Err

	if err != nil {
		logError(err.Error())
		return resp, err
	}

	resp = strings.Split(resp, fmt.Sprintf("%s", commandType))[1]
	resp = strings.Trim(resp, "\r\n")
	resp = strings.TrimSpace(resp)

	return resp, nil
}

// GetActiveSignalByPort checks if the signal on a given port is active or not
func GetActiveSignalByPort(address, port string, readWelcome bool) (structs.ActiveSignal, *nerr.E) {
	var signal structs.ActiveSignal

	signal.Active = false

	signalResponse, err := hardwareCommand(Signal, port, address, readWelcome)
	if err != nil {
		log.L.Error(err.Error())
		return signal, nerr.Translate(err).Addf("failed to get the signal for %s on %s", port, address)
	}

	signalStatus := strings.Split(signalResponse, ",")[1]

	if signalStatus == "1" {
		signal.Active = true
	}

	return signal, nil
}
