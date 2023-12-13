package main

import (
	"github.com/tealeg/xlsx/v3"
)

type SidGroup map[string]struct{}

func parseSids(sh *xlsx.Sheet) []SidGroup {
	groups := make([]SidGroup, 0)
	for g := 0; g < sh.MaxCol; g++ {
		sids := make(map[string]struct{})

		for i := 1; i < sh.MaxRow; i++ {
			c, err := sh.Cell(i, g)
			if err != nil || c.Value == "" {
				break
			}

			sid := c.Value[:len(c.Value)-2]
			sids[sid] = struct{}{}
		}

		groups = append(groups, sids)
	}

	return groups
}
