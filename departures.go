package main

import (
	"fmt"

	"github.com/tealeg/xlsx/v3"
)

type DepartureData struct {
	Callsign string
	ToD      int
	Type     string
	Wake     int    // 0 -> light; 1 -> medium; 2 -> heavy
	SID      string // XXXXX
	Runway   bool   // false -> right; true -> left
}

func parseDepartures(sh *xlsx.Sheet) map[string][]DepartureData {
	departures := make(map[string][]DepartureData, 0)
	err := sh.ForEachRow(func(r *xlsx.Row) error {
		if r.GetCoordinate() == 0 {
			return nil
		}

		callsign := r.GetCell(1).Value
		excelTod, err := r.GetCell(2).Float()
		if err != nil {
			return fmt.Errorf("getting time of departure: %w", err)
		}
		tod := xlsx.TimeFromExcelTime(excelTod, false)
		todSeconds := tod.Hour()*3600 + tod.Minute()*60 + tod.Second()

		typ := r.GetCell(4).Value

		wakeStr := r.GetCell(5).Value
		var wake int
		switch wakeStr {
		case "Ligera":
			wake = 0
		case "Media":
			wake = 1
		case "Pesada":
			wake = 2
		default:
			return fmt.Errorf("unknown wake: %s", wakeStr)
		}

		depProc := r.GetCell(6).Value
		sid := "-"
		if depProc != "-" {
			sid = depProc[:len(depProc)-2]
		}

		var runway bool
		rwStr := r.GetCell(7).Value
		switch rwStr[len(rwStr)-1] {
		case 'R':
			runway = false
		case 'L':
			runway = true
		default:
			return fmt.Errorf("unknown runway: %c", depProc[6])
		}

		depData := DepartureData{
			Callsign: callsign,
			ToD:      todSeconds,
			Type:     typ,
			Wake:     wake,
			SID:      sid,
			Runway:   runway,
		}

		deps, ok := departures[callsign]
		if !ok {
			deps = make([]DepartureData, 0)
		}
		deps = append(deps, depData)
		departures[callsign] = deps

		fmt.Println(depData)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return departures
}
