package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/umahmood/haversine"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

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

//func main() {
//	output := ""
//	//Test Case - 1
//	//Route:
//	route1 := "KSFO +SSTIK4 NTELL Q174 FLCHR COKTL1 KLAS"
//	//Expected Parsed Route:
//	route1expected := "KSFO SSTIK NTELL CABAB TTMSN SKANN FLCHR COKTL KLAS"
//
//	route1expecteddist := 393.9
//	route1dist := Routedist(Routeparse(route1))
//
//	route1output := strings.Join(Routeparse(route1), " ")
//
//	output += verify(1, route1expected, route1output, route1expecteddist, route1dist)
//
//	//Test Case - 2
//	//Route:
//	route2 := "KSEA SEA7 SEA DCT NORMY J90 MWH/N0451F350 DCT KU87M DCT IDA DCT MJANE DCT KDEN"
//	//Expected Parsed Route:
//	route2expected := "KSEA SEA NORMY BLUIT MWH KU87M IDA MJANE KDEN"
//
//	route2expecteddist := 902.8
//	route2dist := Routedist(Routeparse(route2))
//
//	route2output := strings.Join(Routeparse(route2), " ")
//
//	output += verify(2, route2expected, route2output, route2expecteddist, route2dist)
//
//	//Test Case - 3
//	//Route:
//	route3 := "KDFW SWABR8 HULZE TXO J72 ABQ J78 DRK GABBL HLYWD1 KLAX"
//	//Expected Parsed Route:
//	route3expected := "KDFW SWABR HULZE TXO MIERA ABQ ZUN PYRIT DRK GABBL HLYWD KLAX"
//
//	route3expecteddist := 1090.3
//	route3dist := Routedist(Routeparse(route3))
//
//	route3output := strings.Join(Routeparse(route3), " ")
//
//	output += verify(3, route3expected, route3output, route3expecteddist, route3dist)
//
//	//Test Case - 4
//	//Route:
//	route4 := "KSFO TRUKN2 MOGEE Q122 KURSE/N0457F350 Q122 FOD J94 PMM J70 LVZ LENDY6 KJFK"
//	//Expected Parsed Route:
//	route4expected := "KSFO TRUKN MOGEE MACUS MCORD LCU BEARR KURSE ONL FOD VIGGR DBQ COTON OBK KUBBS PMM ALPHE DUNKS SVM CFGFT BEWEL JHW HOXIE DMACK STENT MAGIO LVZ LENDY KJFK"
//
//	route4expecteddist := 2271.7
//	route4dist := Routedist(Routeparse(route4))
//
//	route4output := strings.Join(Routeparse(route4), " ")
//
//	output += verify(4, route4expected, route4output, route4expecteddist, route4dist)
//
//	//Test Case - 5
//	//Route:
//	route5 := "KDEN COORZ6 VOAXA Q136 OAL MOD9 KSFO"
//	//Expected Parsed Route
//	route5expected := "KDEN COORZ VOAXA ELLFF WEEMN MANRD TRALP GDGET CRLES KATTS RUMPS OAL KSFO"
//	// TODO -- ^^ CHECK COORZ6 and MOD9 CHARTS
//
//	route5expecteddist := 843.1
//	route5dist := Routedist(Routeparse(route5))
//
//	route5output := strings.Join(Routeparse(route5), " ")
//
//	output += verify(5, route5expected, route5output, route5expecteddist, route5dist)
//
//	//Test Case - 5
//	//Route:
//	route6 := "CYYZ N0472F400 MIXUT6 GNTRY DCT SVM J70 DUNKS DCT BEJAE DCT PMM DCT KG75K DCT KP72G DCT OBH BRWRY LAWGR3 KDEN"
//	//Expected Parsed Route
//	route6expected := "KDEN COORZ VOAXA ELLFF WEEMN MANRD TRALP GDGET CRLES KATTS RUMPS OAL KSFO"
//	// TODO -- ^^ CHECK COORZ6 and MOD9 CHARTS
//
//	route6expecteddist := 1170.4
//	route6dist := Routedist(Routeparse(route6))
//
//	route6output := strings.Join(Routeparse(route6), " ")
//
//	output += verify(6, route6expected, route6output, route6expecteddist, route6dist)
//
//	f, err := os.Create("output.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer f.Close()
//
//	_, err = f.WriteString(output)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Done!")
//}

//Verifies the test cases
func verify(routeID int, expectedroute string, route string, expecteddist float64, dist float64) string {
	output := ""
	if expectedroute == route {
		output += "Route " + strconv.Itoa(routeID) + " -- Parsed Correctly\n"
	} else {
		output += "Route " + strconv.Itoa(routeID) + " -- ERROR\n"
		output += "\tExpected: " + expectedroute + "\n"
		output += "\tOutput: " + route + "\n"
	}
	if math.Abs(expecteddist-dist) < expecteddist*.01 {
		output += "Route " + strconv.Itoa(routeID) + " -- Distance Calculated Correctly\n\n"
	} else {
		output += "Route " + strconv.Itoa(routeID) + " -- ERROR\n"
		output += fmt.Sprintf("\tExpected: %f", expecteddist)
		output += fmt.Sprintf("\tOutput: %f\n\n", dist)
	}
	return output
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

// Routeparse TODO - Destination repeated multiple times
func Routeparse(route string) []string {
	intersections := strings.Fields(route)
	for i := 0; i < len(intersections); i++ {
		if strings.Index(intersections[i], "/") >= 0 {
			intersections[i] = intersections[i][0:strings.Index(intersections[i], "/")]
		}
		if strings.Index(intersections[i], "+") >= 0 {
			intersections[i] = intersections[i][strings.Index(intersections[i], "+")+1:]
		}
		if len(intersections[i]) > 5 {
			intersections[i] = intersections[i][0:5]
		}
	}
	var endroute []string

	//Need to add AWY conversion parser.
	for i := 0; i < len(intersections); i++ {
		if len(airways[intersections[i]]) != 0 {
			start := intersections[i-1]
			end := intersections[i+1]
			var airway []string
			between := false
			for _, s := range airways[intersections[i]] {
				if (s == start || s == end) && !between {
					between = true
					continue
				}
				if s == end {
					between = false
					break
				}
				if s == start {
					between = false
					for i := 0; i < len(airway)/2; i++ {
						j := len(airway) - i - 1
						airway[i], airway[j] = airway[j], airway[i]
					}
					break
				}
				if between {
					airway = append(airway, s)
				}
			}
			endroute = append(endroute, airway...)
		} else {
			endroute = append(endroute, intersections[i])
		}
	}

	for i, s := range endroute {
		if fixes[s].Lon == 0 {
			endroute = removeIndex(endroute, i)
		}
	}

	return endroute
}

func Routedist(route []string) float64 {
	dist := 0.0
	for i := 0; i < len(route)-1; i++ {
		point1 := haversine.Coord{Lat: fixes[route[i]].Lat, Lon: fixes[route[i]].Lon}
		point2 := haversine.Coord{Lat: fixes[route[i+1]].Lat, Lon: fixes[route[i+1]].Lon}
		mi, _ := haversine.Distance(point1, point2)
		mi /= 1.151
		dist += mi
	}

	return dist
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
