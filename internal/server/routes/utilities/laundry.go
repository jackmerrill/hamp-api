package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type MachineType string

const (
	Washer MachineType = "Washer"
	Dryer  MachineType = "Dryer"
)

type Machine struct {
	Name          string      `json:"name"`
	Type          MachineType `json:"type"`
	Status        string      `json:"status"`
	Time          *string     `json:"time"`
	EstimatedTime *time.Time  `json:"estimatedTime"`
}

// LaundryRoom model info
// @Description The model for a laundry room.
type LaundryRoom struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	NextUpdate time.Time `json:"nextUpdate"`
	LastUpdate time.Time `json:"lastUpdate"`

	Machines []Machine `json:"machines"`

	updateChan chan bool
}

func (l *LaundryRoom) GetMachines() error {
	// create the update channel if it doesn't exist
	if l.updateChan == nil {
		l.updateChan = make(chan bool)
	}

	res, err := http.Get(fmt.Sprintf("https://laundrytrackerconnect.com/hamp%s.aspx", l.ID))

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Laundry room not found")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return err
	}

	machines := []Machine{}

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		// skip the first two rows and the last row
		if i < 2 {
			return
		}

		machine := Machine{}

		// get the machine details (name -> .name, type -> .type, status -> .status, time -> .time, notification -> .form)
		// machines only have a single class, but there are three rows that dont have a class. those aren't machines
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				if s.HasClass("name") {
					machine.Name = s.Text()
				}
			case 1:
				if s.HasClass("type") {
					machine.Type = MachineType(s.Text())
				}
			case 2:
				if s.HasClass("status") {
					machine.Status = s.Text()
				}
			case 3:
				if s.HasClass("time") {
					// remove u00a0 (non-breaking space) from the time
					text := strings.ReplaceAll(s.Text(), "\u00a0", "")

					if text != "" {
						machine.Time = &text

						// parse the time (formatted as "MM minutes left")
						minutesLeft := strings.Split(text, " ")[0]

						// parse the minutes left
						minutes, err := time.ParseDuration(fmt.Sprintf("%sm", minutesLeft))

						if err != nil {
							return
						}

						// set the estimated time
						estTime := time.Now().Add(minutes)

						machine.EstimatedTime = &estTime
					} else {
						machine.Time = nil
					}
				}
			}
		})

		if machine.Name != "" {
			machines = append(machines, machine)
		}
	})

	l.Machines = machines

	// set the last update time
	l.LastUpdate = time.Now()

	// set the next update time
	l.NextUpdate = time.Now().Add(45 * time.Second)

	// send an update to the websocket
	l.updateChan <- true

	return nil
}

func (l *LaundryRoom) GetMachine(machine string) (*Machine, error) {
	for _, m := range l.Machines {
		if m.Name == machine {
			return &m, nil
		}
	}

	return nil, fmt.Errorf("Machine not found")
}

func (l *LaundryRoom) StartFetchLoop() {
	// fetches the laundry room data every 45 seconds
	go func() {
		for {
			l.GetMachines()
			time.Sleep(45 * time.Second)
		}
	}()
}

var LaundryRooms = []LaundryRoom{
	{
		ID:   "dakink",
		Name: "Dakin",
	},
	{
		ID:   "merrill",
		Name: "Merrill",
	},
	{
		ID:   "enfield",
		Name: "Enfield",
	},
	{
		ID:   "prescott",
		Name: "Prescott",
	},
}

func InitLaundry() {
	for i := range LaundryRooms {
		LaundryRooms[i].StartFetchLoop()
	}
}

// GetDakin godoc
// @Summary Get the laundry room data for Dakin.
// @Description Get the laundry room data for Dakin. Add /live for a live websocket stream.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} LaundryRoom
// @Router /utilities/laundry/dakin [get]
func GetDakin(c echo.Context) error {
	// get the laundry room
	room := LaundryRooms[0]

	// check if the cache query param is false
	if c.QueryParam("cache") == "false" {
		// get the latest data
		room.GetMachines()
	}

	return c.JSON(200, room)
}

func GetDakinLive(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		room := LaundryRooms[0]

		// send the current machine data
		err := websocket.JSON.Send(ws, room)

		if err != nil {
			return
		}

		for {
			// wait for the next update
			<-room.updateChan

			// send the updated machine data
			err = websocket.JSON.Send(ws, room)

			if err != nil {
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

// GetDakinMachine godoc
// @Summary Get a specific machine in the Dakin laundry room.
// @Description Get a specific machine in the Dakin laundry room.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} LaundryRoom
// @Param id path string true "Machine ID"
// @Router /utilities/laundry/dakin/{id} [get]
func GetDakinMachine(c echo.Context) error {
	room := LaundryRooms[0]

	// get the machine number
	machineNumber := c.Param("machine")

	// get the machine
	machine, err := room.GetMachine(machineNumber)

	if err != nil {
		return c.String(404, err.Error())
	}

	return c.JSON(200, machine)
}

// GetMerrill godoc
// @Summary Get the laundry room data for Merrill.
// @Description Get the laundry room data for Merrill. Add /live for a live websocket stream.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} Machine
// @Router /utilities/laundry/merrill [get]
func GetMerrill(c echo.Context) error {
	// get the laundry room
	room := LaundryRooms[1]

	// check if the cache query param is false
	if c.QueryParam("cache") == "false" {
		// get the latest data
		room.GetMachines()
	}

	return c.JSON(200, room)
}

func GetMerrillLive(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		room := LaundryRooms[1]

		// send the current machine data
		err := websocket.JSON.Send(ws, room)

		if err != nil {
			return
		}

		for {
			// wait for the next update
			<-room.updateChan

			// send the updated machine data
			err = websocket.JSON.Send(ws, room)

			if err != nil {
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

// GetMerrillMachine godoc
// @Summary Get a specific machine in the Merrill laundry room.
// @Description Get a specific machine in the Merrill laundry room.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} Machine
// @Param id path string true "Machine ID"
// @Router /utilities/laundry/merrill/{id} [get]
func GetMerrillMachine(c echo.Context) error {
	room := LaundryRooms[1]

	// get the machine number
	machineNumber := c.Param("machine")

	// get the machine
	machine, err := room.GetMachine(machineNumber)

	if err != nil {
		return c.String(404, err.Error())
	}

	return c.JSON(200, machine)
}

// GetEnfield godoc
// @Summary Get the laundry room data for Enfield.
// @Description Get the laundry room data for Enfield. Add /live for a live websocket stream.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} LaundryRoom
// @Router /utilities/laundry/enfield [get]
func GetEnfield(c echo.Context) error {
	// get the laundry room
	room := LaundryRooms[2]

	// check if the cache query param is false
	if c.QueryParam("cache") == "false" {
		// get the latest data
		room.GetMachines()
	}

	return c.JSON(200, room)
}

func GetEnfieldLive(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		room := LaundryRooms[2]

		// send the current machine data
		err := websocket.JSON.Send(ws, room)

		if err != nil {
			return
		}

		for {
			// wait for the next update
			<-room.updateChan

			// send the updated machine data
			err = websocket.JSON.Send(ws, room)

			if err != nil {
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

// GetEnfieldMachine godoc
// @Summary Get a specific machine in the Enfield laundry room.
// @Description Get a specific machine in the Enfield laundry room.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} Machine
// @Param id path string true "Machine ID"
// @Router /utilities/laundry/enfield/{id} [get]
func GetEnfieldMachine(c echo.Context) error {
	room := LaundryRooms[2]

	// get the machine number
	machineNumber := c.Param("machine")

	// get the machine
	machine, err := room.GetMachine(machineNumber)

	if err != nil {
		return c.String(404, err.Error())
	}

	return c.JSON(200, machine)
}

// GetPrescott godoc
// @Summary Get the laundry room data for Prescott.
// @Description Get the laundry room data for Prescott. Add /live for a live websocket stream.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} LaundryRoom
// @Router /utilities/laundry/prescott [get]
func GetPrescott(c echo.Context) error {
	// get the laundry room
	room := LaundryRooms[3]

	// check if the cache query param is false
	if c.QueryParam("cache") == "false" {
		// get the latest data
		room.GetMachines()
	}

	return c.JSON(200, room)
}

func GetPrescottLive(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		room := LaundryRooms[3]

		// send the current machine data
		err := websocket.JSON.Send(ws, room)

		if err != nil {
			return
		}

		for {
			// wait for the next update
			<-room.updateChan

			// send the updated machine data
			err = websocket.JSON.Send(ws, room)

			if err != nil {
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

// GetPrescottMachine godoc
// @Summary Get a specific machine in the Prescott laundry room.
// @Description Get a specific machine in the Prescott laundry room.
// @Tags utilities
// @Accept json
// @Produce json
// @Success 200 {object} Machine
// @Param id path string true "Machine ID"
// @Router /utilities/laundry/prescott/{id} [get]
func GetPrescottMachine(c echo.Context) error {
	room := LaundryRooms[3]

	// get the machine number
	machineNumber := c.Param("machine")

	// get the machine
	machine, err := room.GetMachine(machineNumber)

	if err != nil {
		return c.String(404, err.Error())
	}

	return c.JSON(200, machine)
}
