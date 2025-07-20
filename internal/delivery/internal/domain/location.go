package domain

import "time"

type Location struct {
	Id        uint
	Latitude  float64
	Longitude float64
	Timestamp int64
}

func NewLocation(latitude float64, longitude float64) *Location {
	return &Location{
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now().UnixMilli(),
	}
}
