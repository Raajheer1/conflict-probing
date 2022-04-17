package main

import (
	"fmt"
	"github.com/umahmood/haversine"
	"math"
)

func points(flp FlightPlan) []Location {
	if flp.Points == nil {
		var points []Location
		for i := 0; i < len(flp.RteParse)-1; i++ {
			point1 := haversine.Coord{Lat: fixes[flp.RteParse[i]].Lat, Lon: fixes[flp.RteParse[i]].Lon}
			point2 := haversine.Coord{Lat: fixes[flp.RteParse[i+1]].Lat, Lon: fixes[flp.RteParse[i+1]].Lon}
			mi, _ := haversine.Distance(point1, point2)
			mi /= 1.151
			dist := mi

			for x := 0; x < int(dist); x++ {
				points = append(points, Location{
					math.Abs(fixes[flp.RteParse[i]].Lat-fixes[flp.RteParse[i+1]].Lat)/10.0*float64(x) + fixes[flp.RteParse[i+1]].Lon,
					math.Abs(fixes[flp.RteParse[i]].Lon-fixes[flp.RteParse[i+1]].Lon)/float64(int(dist))*float64(x) + fixes[flp.RteParse[i+1]].Lon,
				})
			}
		}
		//fmt.Println(points)
		return points
	}
	fmt.Println("PRINT THIS")
	return []Location{}
}
