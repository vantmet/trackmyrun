package main

import "strconv"

type Run struct {
	Distance int //All distances stored in km.
}

func GetRunDistanceKm(r Run) string {
	return strconv.Itoa(r.Distance) + "km"
}
