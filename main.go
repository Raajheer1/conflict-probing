package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FlightPlan struct {
	Rules     string `json:"flight_rules"`
	Aircraft  string `json:"aircraft_faa"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Alternate string `json:"alternate"`
	TAS       string `json:"cruise_tas"`
	Altitude  string `json:"altitude"`
	Alt       int    `json:"alt"`
	Deptime   string `json:"deptime"`
	ERT       string `json:"enroute_time"`
	Fuel      string `json:"fuel_time"`
	Remarks   string `json:"remarks"`
	Route     string `json:"route"`
}

type Aircraft struct {
	CID        uint       `json:"cid"`
	Callsign   string     `json:"callsign"`
	Latitude   float32    `json:"latitude"`
	Longitude  float32    `json:"longitude"`
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

type Location struct {
	Lon float64 `xml:"Lon,attr"`
	Lat float64 `xml:"Lat,attr"`
}

type Fix struct {
	Name     string   `xml:"ID,attr"`
	Location Location `xml:"Location"`
}

type Fixes struct {
	Fixes []Fix `xml:"Waypoint"`
}

type Airway struct {
	Name  string
	Fixes []string
}

var fixes map[string]Location = parseFIX("Waypoints.xml")
var airways map[string][]string = parseAWY("AWY.txt")

func main() {
	fmt.Println("Initializing...")
	start := time.Now()
	//Grab initial aircraft data.
	var aircraft Aircrafts = initialize()

	//Wait for next cycle
	time.Sleep(20 * time.Second)
	//GRAB next cycle

	//Compare Aircraft / Check for changes

	duration := time.Since(start)
	fmt.Print("Runtime: ")
	fmt.Println(duration)

}

func initialize() Aircrafts {
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

		for i := 0; i < len(aircraft.Aircrafts); i++ {
			if aircraft.Aircrafts[i].Flightplan.Rules == "I" {
				altConversion(&aircraft.Aircrafts[i])
			}
		}

		//for i := 0; i < len(aircraft.Aircrafts); i++ {
		//	if aircraft.Aircrafts[i].Flightplan.Route != "" && aircraft.Aircrafts[i].Flightplan.Rules == "I" && aircraft.Aircrafts[i].Flightplan.Arrival == "KDEN" {
		//		fmt.Println("Callsign: " + aircraft.Aircrafts[i].Callsign)
		//		fmt.Println("Route: " + aircraft.Aircrafts[i].Flightplan.Departure + " " + aircraft.Aircrafts[i].Flightplan.Route + " " + aircraft.Aircrafts[i].Flightplan.Arrival)
		//		aircraft.Aircrafts[i].Distance = Routedist(Routeparse(aircraft.Aircrafts[i].Flightplan.Departure + " " + aircraft.Aircrafts[i].Flightplan.Route + " " + aircraft.Aircrafts[i].Flightplan.Arrival))
		//		fmt.Printf("Distance: %.2f", aircraft.Aircrafts[i].Distance)
		//		fmt.Println("\n---")
		//	}
		//}
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

//Parses the latest AIRAC FIXES
func parseFIX(filename string) map[string]Location {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var fixesXML Fixes
	xml.Unmarshal(byteValue, &fixesXML)

	fixesmap := make(map[string]Location)

	for i := 0; i < len(fixesXML.Fixes); i++ {
		name := fixesXML.Fixes[i].Name
		loc := fixesXML.Fixes[i].Location
		fixesmap[name] = loc
	}

	return fixesmap
}

//Parses the latest AIRAC AWYs
func parseAWY(filename string) map[string][]string {
	txtFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	defer txtFile.Close()
	scanner := bufio.NewScanner(txtFile)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	airwaysmap := make(map[string][]string)
	for _, eachline := range txtlines {
		name := eachline[1:strings.Index(eachline, "F")]
		fixeslist := strings.Fields(eachline[strings.Index(eachline, "F")+6:])

		airwaysmap[name] = fixeslist
	}

	return airwaysmap
}
