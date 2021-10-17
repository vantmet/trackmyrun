package main

import "strconv"

type RunTime struct {
	Hours   int
	Minutes int
	Seconds float32
}

type Run struct {
	Distance int //All distances stored in km.
	RunTime  RunTime
}

func GetRunDistanceKm(r Run) string {
	return strconv.Itoa(r.Distance) + "km"
}

func GetRunTime(r Run) RunTime {
	return r.RunTime
}
