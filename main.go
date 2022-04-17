package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FlightPlan struct {
	Rules     string     `json:"flight_rules"`
	Aircraft  string     `json:"aircraft_faa"`
	Departure string     `json:"departure"`
	Arrival   string     `json:"arrival"`
	Alternate string     `json:"alternate"`
	TAS       string     `json:"cruise_tas"`
	Altitude  string     `json:"altitude"`
	Alt       int        `json:"alt"`
	Deptime   string     `json:"deptime"`
	ERT       string     `json:"enroute_time"`
	Fuel      string     `json:"fuel_time"`
	Remarks   string     `json:"remarks"`
	Route     string     `json:"route"`
	RteParse  []string   `json:"parsed_route"`
	Points    []Location `json:"points"`
	DCT       string     `json:"direct"`
}

type Aircraft struct {
	CID        uint       `json:"cid"`
	Callsign   string     `json:"callsign"`
	Latitude   float32    `json:"latitude"`
	OldLat     float32    `json:"OldLat"`
	Longitude  float32    `json:"longitude"`
	OldLon     float32    `json:"OldLon"`
	Altitude   int        `json:"altitude"`
	Groudspeed uint       `json:"groudspeed"`
	Heading    uint       `json:"heading"`
	Flightplan FlightPlan `json:"flight_plan"`
	Status     string     `json:"status"`
	EDT        uint       `json:"edt"`
	DT         uint       `json:"dt"`
	Distance   float64    `json:"distance"`
	Arrival    uint       `json:"arrival"`
}

type Aircrafts struct {
	Aircrafts []Aircraft `json:"pilots"`
}

func main() {
	fmt.Println("Initializing...")
	start := time.Now()
	//Grab initial aircraft data.
	var aircraft Aircrafts = fetchPlanes()
	for _, plane := range aircraft.Aircrafts {
		plane.Flightplan.RteParse = Routeparse(plane.Flightplan.Departure + " " + plane.Flightplan.Route + " " + plane.Flightplan.Arrival)
	}

	//Wait for next cycle
	time.Sleep(20 * time.Second)
	var newAircraft Aircrafts = fetchPlanes()
	//GRAB next cycle

	//Compare Aircraft / Check for changes
	current := differences(aircraft.Aircrafts, newAircraft.Aircrafts)

	for i := 0; i < len(current); i++ {
		current[i].Flightplan.Points = points(current[i].Flightplan)
	}

	output, _ := json.MarshalIndent(current, "", " ")
	file, _ := os.Create("output.json")
	_, _ = file.WriteString(string(output))

	duration := time.Since(start)
	fmt.Print("Runtime: ")
	fmt.Println(duration)

}

func fetchPlanes() Aircrafts {
	resp, err := http.Get("https://data.vatsim.net/v3/vatsim-data.json")
	if err != nil {
		fmt.Printf("HTTP GET error: ", err)
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}

		var aircraft Aircrafts
		err = json.Unmarshal(contents, &aircraft)
		if err != nil {
			fmt.Println(err)
		}

		//Filter out aircraft that are VFR
		for i := 0; i < len(aircraft.Aircrafts); i++ {
			fmt.Println(aircraft.Aircrafts[i].Callsign)
			if aircraft.Aircrafts[i].Flightplan.Rules != "I" {
				if len(aircraft.Aircrafts) > 2 {
					aircraft.Aircrafts = append(aircraft.Aircrafts[:i], aircraft.Aircrafts[i+1:]...)
					i--
				} else {
					if i == 1 {
						aircraft.Aircrafts = aircraft.Aircrafts[:1]
					} else {
						aircraft.Aircrafts = aircraft.Aircrafts[2:]
					}
					i--
				}
				continue
			}
			if aircraft.Aircrafts[i].Flightplan.Rules == "I" {
				altConversion(&aircraft.Aircrafts[i])
			}

			//This is causing issues for the last few aircraft
			if aircraft.Aircrafts[i].Flightplan.Departure == "" || aircraft.Aircrafts[i].Flightplan.Arrival == "" {
				if len(aircraft.Aircrafts) > 2 {
					aircraft.Aircrafts = append(aircraft.Aircrafts[:i], aircraft.Aircrafts[i+1:]...)
					i--
				} else {
					if i == 1 {
						aircraft.Aircrafts = aircraft.Aircrafts[:1]
					} else {
						aircraft.Aircrafts = aircraft.Aircrafts[2:]
					}
					i--
				}
			} else {
				if aircraft.Aircrafts[i].Flightplan.Departure[:1] != "K" {
					if len(aircraft.Aircrafts) > 2 {
						aircraft.Aircrafts = append(aircraft.Aircrafts[:i], aircraft.Aircrafts[i+1:]...)
						i--
					} else {
						if i == 1 {
							aircraft.Aircrafts = aircraft.Aircrafts[:1]
						} else {
							aircraft.Aircrafts = aircraft.Aircrafts[2:]
						}
						i--
					}
				} else if aircraft.Aircrafts[i].Flightplan.Arrival[:1] != "K" {
					if len(aircraft.Aircrafts) > 2 {
						aircraft.Aircrafts = append(aircraft.Aircrafts[:i], aircraft.Aircrafts[i+1:]...)
						i--
					} else {
						if i == 1 {
							aircraft.Aircrafts = aircraft.Aircrafts[:1]
						} else {
							aircraft.Aircrafts = aircraft.Aircrafts[2:]
						}
						i--
					}
				}
			}
		}

		//TODO- Add filter for aircraft within the USA

		flightStatus(&aircraft)
		return aircraft
	}

	return Aircrafts{}
}

//Converts altitude string to Int String
func altConversion(aircraft *Aircraft) {
	if strings.Contains(aircraft.Flightplan.Altitude, "FL") {
		altitude, err := strconv.Atoi(aircraft.Flightplan.Altitude[2:len(aircraft.Flightplan.Altitude)])
		if err != nil {
			fmt.Println("Error parsing Altitude", err)
		}
		aircraft.Flightplan.Alt = altitude * 100
	} else {
		altitude, err := strconv.Atoi(aircraft.Flightplan.Altitude)
		if altitude < 1000 {
			altitude *= 100
		}
		if err != nil {
			fmt.Println("Error parsing Altitude", err)
		}
		aircraft.Flightplan.Alt = altitude
	}
}

//Checks to see if at cruise
func flightStatus(aircraft *Aircrafts) {
	for i := 0; i < len(aircraft.Aircrafts); i++ {
		if aircraft.Aircrafts[i].Flightplan.Rules == "I" {
			if aircraft.Aircrafts[i].Altitude > (aircraft.Aircrafts[i].Flightplan.Alt - aircraft.Aircrafts[i].Flightplan.Alt/100) {
				//fmt.Println(aircraft.Aircrafts[i].Altitude, " --- ", aircraft.Aircrafts[i].Flightplan.Alt, " --- ", aircraft.Aircrafts[i].Flightplan.Alt-aircraft.Aircrafts[i].Flightplan.Alt/100)
				aircraft.Aircrafts[i].Status = "cruise"
			}
		}
	}
}
