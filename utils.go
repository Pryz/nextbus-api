package main

import (
	"time"
)


//TODO: Move to nextbus API
const WEEKDAY string = "wkd"
const SATURDAY string = "sat"
const SUNDAY string = "sun"


// Return Nextbus ServiceClass from specified week day
func dayToServiceClass(day string) string {
	week := [5]string{"monday", "tuesday", "wednesday", "thursday", "friday"}
	for _, d := range week {
		if day == d {
			return WEEKDAY
		}
	}
	if day == "saturday" {
		return SATURDAY
	}
	if day == "sunday" {
		return SUNDAY
	}
	return ""
}


// Generate Epoch time from a string
// Example : 1:18AM -> 4080000
func stringToEpoch(s string) int {
	t0, _ := time.Parse(time.Kitchen, "0:00AM")
	t, _ := time.Parse(time.Kitchen, s)
	return int(t.Sub(t0).Seconds() * 1000)
}

