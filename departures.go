package main

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

type DepartureData struct {
	Callsign string
	ToD      int
	Type     string
	Wake     int    // 0 -> light; 1 -> medium; 2 -> heavy
	SID      string // XXXXX
	Runway   string
}

func parseDepartures(sh *xlsx.Sheet, sids map[string]struct{}) []DepartureData {
	departures := make([]DepartureData, 0)
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
		} else {
			words := strings.Split(r.GetCell(3).Value, " ")
			for _, word := range words {
				if strings.Contains(word, "(") && strings.Contains(word, ")") {
					for i := 0; i < len(word)-2; i++ {
						if word[i] == '(' {
							word = word[i+1:]
							break
						}
					}

					for i := len(word) - 1; i > 0; i-- {
						if word[i] == ')' {
							word = word[:i]
							break
						}
					}

				}

				if _, ok := sids[word]; ok {
					sid = word
					break
				}
			}
		}

		if sid == "-" {
			panic("Couldn't associate a SID to a departure")
		}

		runway := r.GetCell(7).Value[5:]

		if runway != "24L" && runway != "6R" {
			return nil
		}

		depData := DepartureData{
			Callsign: callsign,
			ToD:      todSeconds,
			Type:     typ,
			Wake:     wake,
			SID:      sid,
			Runway:   runway,
		}

		departures = append(departures, depData)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return departures
}
