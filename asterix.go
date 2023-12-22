package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type AsterixData struct {
	GPSCoords
	Time float64

	Callsign string
}

func readAsterix(file *os.File) []AsterixData {
	reader := csv.NewReader(file)
	reader.Comma = ';'
	_, err := reader.Read()
	if err != nil {
		panic(err)
	}
	asterixData := make([]AsterixData, 0)
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		lat, _ := strconv.ParseFloat(strings.Replace(row[0], ",", ".", 1), 64)
		lon, _ := strconv.ParseFloat(strings.Replace(row[1], ",", ".", 1), 64)
		alt, _ := strconv.ParseFloat(strings.Replace(row[2], ",", ".", 1), 64)
		time, _ := strconv.ParseFloat(strings.Replace(row[3], ",", ".", 1), 64)
		callsign := row[7]

		asterixData = append(asterixData, AsterixData{
			GPSCoords: GPSCoords{
				Lat: lat,
				Lon: lon,
				Alt: alt,
			},
			Time:     time,
			Callsign: callsign,
		})
	}

	return asterixData
}
