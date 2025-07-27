package color

import (
	"fmt"
	"math"
	"strings"
)

type Color struct {
	Name string
	Hex  string
	R    int
	G    int
	B    int
}

func HexToRGB(hex string) (int, int, int, error) {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b int
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b, err
}

func ClosestColorName(r, g, b int, cssColors []Color) string {
	closestName := ""
	minDistance := math.MaxFloat64

	for _, color := range cssColors {
		distance := math.Sqrt(float64((r-color.R)*(r-color.R) + (g-color.G)*(g-color.G) + (b-color.B)*(b-color.B)))
		if distance < minDistance {
			minDistance = distance
			closestName = color.Name
		}
	}
	return closestName
}
