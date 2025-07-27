package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelvinzer0/palette-mcp/internal/color"
	"github.com/kelvinzer0/palette-mcp/internal/data"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: colorname #hexcode")
		os.Exit(1)
	}

	hexCode := os.Args[1]
	r, g, b, err := color.HexToRGB(hexCode)
	if err != nil {
		fmt.Println("Invalid hex code:", hexCode)
		os.Exit(1)
	}

	colorName := color.ClosestColorName(r, g, b, data.CssColors)

	// Output JSON for system integration
	result := map[string]string{
		"hex":   hexCode,
		"name":  colorName,
		"rgb":   fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
	}
	json.NewEncoder(os.Stdout).Encode(result)
}
