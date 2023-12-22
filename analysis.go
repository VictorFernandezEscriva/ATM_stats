package main

import (
	"math"
)

func getDistances(projection SystemCartesian, flight1 []AsterixData, flight2 []AsterixData) []float64 {
	i := 0
	j := 0

	var distances []float64
	for {
		if i >= len(flight1) || j >= len(flight2) {
			break
		}

		p1, p2 := flight1[i], flight2[j]

		// Check that timestamps are coordinated
		if math.Abs(p1.Time-p2.Time) > 2 {
			if p1.Time < p2.Time {
				i++
			} else {
				j++
			}
			continue
		}

		i++
		j++
		proj1 := projection.GeocentricToStereographic(p1.ToGeocentric())
		proj2 := projection.GeocentricToStereographic(p2.ToGeocentric())

		du := proj2.U - proj1.U
		dv := proj2.V - proj1.V
		distance := math.Sqrt(du*du+dv*dv) / 1852
		distances = append(distances, distance)
	}

	return distances
}
