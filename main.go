package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/tealeg/xlsx/v3"
)

const (
	FP_EXCEL       = "data/2305_02_dep_lebl.xlsx"
	DECODED_CSV    = "data/230205_08_12_decodif_P3.csv"
	AC_CLASS_EXCEL = "data/Tabla_Clasificacion_aeronaves.xlsx"
	SID_EXCEL_R    = "data/Tabla_misma_SID_06R.xlsx"
	SID_EXCEL_L    = "data/Tabla_misma_SID_24L.xlsx"
)

/*
Analizar:
- distancia mínima entre cada 2 despegues consecutivos (y tiempo en el que se produce)
- distancia media para cada
*/

type FlightInfo struct {
	Callsign string
	Company  string
	Wake     string
	Class    string
	SidGroup string
}

type Result struct {
	First       FlightInfo
	Second      FlightInfo
	MinDistance float64
	Separation  float64
}

type ResultsJSON struct {
	Results []Result
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

	sh := wbSidR.Sheets[0]
	sidGroups := parseSids(sh)

	sids := make(map[string]struct{}, 0)
	for _, group := range sidGroups {
		for sid := range group {
			sids[sid] = struct{}{}
		}
	}

	sh = wbFlightplans.Sheets[0]
	departures := parseDepartures(sh, sids)

	fmt.Println("Number of aircraft with departures:", len(departures))

	sh = wbClasses.Sheets[0]
	classes := parseClasses(sh)

	csvFiles := []string{
		"data/csv 0-4.csv",
		"data/csv 4-8.csv",
		"data/csv 8-12.csv",
		"data/csv 12-16.csv",
		"data/csv 16-20.csv",
		"data/csv 20-24.csv",
	}
	var asterixData []AsterixData
	for _, name := range csvFiles {
		// Read Asterix
		file, err := os.Open(name)
		if err != nil {
			log.Fatal("reading decoded csv:", err)
		}

		asterixData = append(asterixData, readAsterix(file)...)
	}
	asterixData = filterData(asterixData, departures)

	sort.Slice(asterixData, func(i, j int) bool {
		return asterixData[i].Time < asterixData[j].Time
	})

	captureStart := asterixData[0].Time
	captureEnd := asterixData[len(asterixData)-1].Time
	fmt.Println(captureStart, captureEnd, len(asterixData))

	projection := NewSystemCartesian(GPSCoords{
		41.0656560,
		1.413301,
		3438.954,
	})

	proj1 := projection.GeocentricToStereographic(GPSCoords{
		41.0656560,
		1.413301,
		3438.954,
	}.ToGeocentric())
	if proj1.U != 0 || proj1.V != 0 {
		log.Fatal("Error with stereographic projection")
	}

	results := ResultsJSON{}
	for i, dep1 := range departures[:len(departures)-1] {
		flight1 := make([]AsterixData, 0)
		// Only analyze flights that we have data from the start
		if dep1.ToD < int(captureStart) || dep1.ToD > int(captureEnd) {
			continue
		}
		//pairsToCheck = append(pairsToCheck, [2]int{i, i + 1})

		for _, d := range asterixData {
			if d.Callsign == dep1.Callsign && d.Time >= float64(dep1.ToD) {
				// Prevent data from another flight
				if len(flight1) != 0 && d.Time-flight1[len(flight1)-1].Time > 3600 {
					break
				}
				flight1 = append(flight1, d)
			}
		}

		if len(flight1) == 0 {
			continue
		}

		//fmt.Printf("%s %.2f\n", dep1.Callsign, flight1[len(flight1)-1].Time-flight1[0].Time)

		dep2 := departures[i+1]
		if dep2.ToD < int(captureStart) || dep2.ToD > int(captureEnd) {
			continue
		}

		flight2 := make([]AsterixData, 0)
		for _, d := range asterixData {
			if d.Callsign == dep2.Callsign && d.Time >= float64(dep2.ToD) {
				// Prevent data from another flight
				if len(flight2) != 0 && d.Time-flight2[len(flight2)-1].Time > 3600 {
					break
				}
				flight2 = append(flight2, d)
			}
		}

		distances := getDistances(projection, flight1, flight2)
		if distances == nil {
			// Planes doesn't coexist in time
			continue
		}

		minDistance := distances[0]
		for _, d := range distances {
			if d < minDistance {
				minDistance = d
			}
		}

		var sid1 string
		for i, g := range sidGroups {
			_, ok := g[dep1.SID]
			if ok {
				sid1 = fmt.Sprintf("G%d", i+1)
				break
			}
		}
		if sid1 == "" {
			log.Fatal("Unknown group for sid:", dep1.SID)
		}
		var sid2 string
		for i, g := range sidGroups {
			_, ok := g[dep2.SID]
			if ok {
				sid2 = fmt.Sprintf("G%d", i+1)
				break
			}
		}
		if sid2 == "" {
			log.Fatal("Unknown group for sid:", dep2.SID)
		}

		var class1 string
		for _, c := range classes {
			_, ok := c.Types[dep1.Type]
			if ok {
				class1 = c.Name
				break
			}
		}
		if class1 == "" {
			class1 = "R"
		}

		var class2 string
		for _, c := range classes {
			_, ok := c.Types[dep2.Type]
			if ok {
				class2 = c.Name
				break
			}
		}
		if class2 == "" {
			class2 = "R"
		}

		results.Results = append(results.Results, Result{
			First: FlightInfo{
				Callsign: dep1.Callsign,
				Company:  dep1.Callsign[:3],
				Wake:     dep1.Wake,
				Class:    class1,
				SidGroup: sid1,
			},
			Second: FlightInfo{
				Callsign: dep2.Callsign,
				Company:  dep2.Callsign[:3],
				Wake:     dep2.Wake,
				Class:    class2,
				SidGroup: sid2,
			},
			MinDistance: minDistance,
			Separation:  distances[0],
		})
	}

	// Save the results to a JSON file
	saveResults(results)
}

func saveResults(results ResultsJSON) {
	jsonResults, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("results.json", jsonResults, 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("python", "plot.py", "results.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing Python script:", err)
		return
	}

	fmt.Println("Python script executed successfully")
}
