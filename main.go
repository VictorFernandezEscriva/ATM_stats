package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

const (
	FP_EXCEL       = "data/2305_02_dep_lebl.xlsx"
	DECODED_CSV    = "data/230205_08_12_decodif_P3.csv"
	AC_CLASS_EXCEL = "data/Tabla_Clasificacion_aeronaves.xlsx"
	SID_EXCEL_R    = "data/Tabla_misma_SID_06R.xlsx"
	SID_EXCEL_L    = "data/Tabla_misma_SID_24L.xlsx"
)

type AsterixData struct {
	Lat  float64
	Lon  float64
	Alt  float64
	Time float64

	Callsign string
}

func main() {
	wbFlightplans, err := xlsx.OpenFile(FP_EXCEL)
	if err != nil {
		log.Fatal("reading flightplans excel:", err)
	}

	wbClasses, err := xlsx.OpenFile(AC_CLASS_EXCEL)
	if err != nil {
		log.Fatal("reading aircraft classes excel:", err)
	}

	wbSidR, err := xlsx.OpenFile(SID_EXCEL_R)
	if err != nil {
		log.Fatal("reading SID R excel:", err)
	}

	_, err = xlsx.OpenFile(SID_EXCEL_L)
	if err != nil {
		log.Fatal("reading SID L excel:", err)
	}

	sh := wbFlightplans.Sheets[0]
	departures := parseDepartures(sh)

	fmt.Println("Number of aircraft with departures:", len(departures))

	sh = wbSidR.Sheets[0]
	//sidGroups := parseSids(sh)
	parseSids(sh)

	sh = wbClasses.Sheets[0]
	//classes := parseClasses(sh)
	parseClasses(sh)

	// Decode CSV
	file, err := os.Open(DECODED_CSV)
	if err != nil {
		log.Fatal("reading decoded csv:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	_, err = reader.Read()
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
			Lat:      lat,
			Lon:      lon,
			Alt:      alt,
			Time:     time,
			Callsign: callsign,
		})
	}
}
