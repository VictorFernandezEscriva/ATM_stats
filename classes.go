package main

import (
	"github.com/tealeg/xlsx/v3"
)

type AircraftClass struct {
	Name  string
	Types map[string]struct{}
}

func parseClasses(sh *xlsx.Sheet) []AircraftClass {
	classes := make([]AircraftClass, 0)
	for g := 0; g < sh.MaxCol; g++ {
		types := make(map[string]struct{})

		c, err := sh.Cell(0, g)
		if err != nil {
			panic(err)
		}
		name := c.Value

		for i := 1; i < sh.MaxRow; i++ {
			c, err := sh.Cell(i, g)
			if err != nil || c.Value == "" {
				break
			}

			t := c.Value
			types[t] = struct{}{}
		}

		classes = append(classes, AircraftClass{
			Name:  name,
			Types: types,
		})
	}

	return classes
}
