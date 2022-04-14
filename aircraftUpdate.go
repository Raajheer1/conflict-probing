package main

import (
	"fmt"
)

//todo - Checks differences between previous cycle and current cycle, returns differences

func Differences(previous []Aircraft, current []Aircraft) {
	fmt.Println("Checking differences...")
	fmt.Print("Current List: ")
	fmt.Println(len(current))
	fmt.Print("Old List: ")
	fmt.Println(len(previous))
	fmt.Println("-----------")
	var combined []Aircraft
	for _, aircraft := range previous {
		for i2, newaircraft := range current {

			//Aircraft still online
			if aircraft.Callsign == newaircraft.Callsign {
				//Check for differences
				temp := newaircraft
				temp.OldLat = aircraft.Latitude
				temp.OldLon = aircraft.Longitude
				combined = append(combined, temp)

				//I dont think this 1st one does anything
				//previous = removeIndex(previous, i1)

				//This def does something
				current = removeIndex(current, i2)
			}
		}
	}

	// At this point the Current List contains only NEW planes

	fmt.Print("Current List: ")
	fmt.Println(len(current))
	fmt.Print("Old List: ")
	fmt.Println(len(previous))
	fmt.Print("Combined List: ")
	fmt.Println(len(combined))
	fmt.Println("-----------")

	for _, aircraft := range combined {
		for i2, aircraft2 := range previous {
			if aircraft.Callsign == aircraft2.Callsign {
				previous = removeIndex(previous, i2)
			}
		}
	}

	// At this point the Previous List should contain only OLD planes who have signed off
	fmt.Print("Current List: ")
	fmt.Println(len(current))
	fmt.Print("Old List: ")
	fmt.Println(len(previous))
	fmt.Print("Combined List: ")
	fmt.Println(len(combined))
	fmt.Println("-----------")

	//fmt.Println("Disconnected List")
	//fmt.Println("----------------------------------------------------")
	//for _, aircraft := range previous {
	//	fmt.Println("- ", aircraft.Callsign)
	//}
	//
	//fmt.Println("New Aircraft List")
	//fmt.Println("----------------------------------------------------")
	//for _, aircraft := range current {
	//	fmt.Println("- ", aircraft.Callsign)
	//}
}

func removeIndex(s []Aircraft, index int) []Aircraft {
	return append(s[:index], s[index+1:]...)
}
