package monitor

import (
	"strconv"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/kramer-microservice/via"
)

//QueryPresentationNumber .
func QueryPresentationNumber(baseEvent events.Event, f func(events.Event)) {

	c, err := via.GetPresenterCount("ITB-1101-VIA1.byu.edu")
	if err != nil {
		log.L.Errorf("Couldn't get presenter count: %v", err.Error())
	}

	baseEvent.Key = "presenter-count"
	baseEvent.Value = strconv.Itoa(c)
	baseEvent.Timestamp = time.Now()

	f(baseEvent)
}
