package main

func filterData(data []AsterixData, departures []DepartureData) []AsterixData {
	departedLEBL := make(map[string]struct{})
	for _, d := range departures {
		departedLEBL[d.Callsign] = struct{}{}
	}

	filtered := make([]AsterixData, 0)
	for _, d := range data {
		// Keep only data from planes that departed from LEBL
		if _, ok := departedLEBL[d.Callsign]; !ok {
			continue
		}

		filtered = append(filtered, d)
	}

	return filtered
}
