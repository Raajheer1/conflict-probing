package main

import (
	"fmt"
)

//todo - Checks differences between previous cycle and current cycle, returns differences

func differences(previous []Aircraft, current []Aircraft) []Aircraft {
	fmt.Println("Checking differences...")
	var combined []Aircraft
	for _, aircraft := range previous {
		for i2, newaircraft := range current {

			//Aircraft still online
			if aircraft.Callsign == newaircraft.Callsign {
				//Check for differences

				//Checks to see if routes have not changed, if no change then copy over parsed route
				if aircraft.Flightplan.Route == newaircraft.Flightplan.Route {
					newaircraft.Flightplan.RteParse = aircraft.Flightplan.RteParse
				} else {
					//Routes have changed, reparse the route
					newaircraft.Flightplan.RteParse = Routeparse(newaircraft.Flightplan.Departure + " " + newaircraft.Flightplan.Route + " " + newaircraft.Flightplan.Arrival)
				}

				temp := newaircraft
				temp.OldLat = aircraft.Latitude
				temp.OldLon = aircraft.Longitude

				combined = append(combined, temp)

				//I dont think this 1st one does anything
				//previous = removeIndex(previous, i1)

				//This def does something
				current = removeAircraftIndex(current, i2)
			}
		}
	}

	//for _, aircraft := range combined {
	//	for i2, aircraft2 := range previous {
	//		if aircraft.Callsign == aircraft2.Callsign {
	//			previous = removeIndex(previous, i2)
	//		}
	//	}
	//}
	//
	//// At this point the Previous List should contain only OLD planes who have signed off
	//fmt.Print("Current List: ")
	//fmt.Println(len(current))
	//fmt.Print("Old List: ")
	//fmt.Println(len(previous))
	//fmt.Print("Combined List: ")
	//fmt.Println(len(combined))
	//fmt.Println("-----------")

	for _, aircraft := range current {
		aircraft.Flightplan.RteParse = Routeparse(aircraft.Flightplan.Departure + " " + aircraft.Flightplan.Route + " " + aircraft.Flightplan.Arrival)
		combined = append(combined, aircraft)
	}

	return combined
}

func removeAircraftIndex(s []Aircraft, index int) []Aircraft {
	return append(s[:index], s[index+1:]...)
}
