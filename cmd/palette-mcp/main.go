package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/kelvinzer0/palette-mcp/internal/color"
	"github.com/kelvinzer0/palette-mcp/internal/data"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			startMcpServer(false, "")
		case "serve-http":
			port := "8080"
			if len(os.Args) > 2 {
				port = os.Args[2]
			}
			startMcpServer(true, port)
		default:
			runCliTool()
		}
	} else {
		runCliTool()
	}
}

func runCliTool() {
	if len(os.Args) < 2 {
		logrus.Fatal("Usage: palette-mcp #hexcode")
	}

	hexCode := os.Args[1]
	r, g, b, err := color.HexToRGB(hexCode)
	if err != nil {
		logrus.Fatal("Invalid hex code:", hexCode)
	}

	colorName := color.ClosestColorName(r, g, b, data.CssColors)

	// Output JSON for system integration
	result := map[string]string{
		"hex":  hexCode,
		"name": colorName,
		"rgb":  fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
	}
	json.NewEncoder(os.Stdout).Encode(result)
}

func startMcpServer(useHTTP bool, port string) {
	s := server.NewMCPServer(
		"Palette MCP",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	getColorInfoTool := mcp.NewTool("get_color_info",
		mcp.WithDescription("Retrieves information about a color given its hex code."),
		mcp.WithString("hexCode",
			mcp.Required(),
			mcp.Description("The hex code of the color (e.g., #FF0000)."),
		),
	)

	s.AddTool(getColorInfoTool, handleGetColorInfo)

	if useHTTP {
		if err := server.NewStreamableHTTPServer(s).Start(":" + port); err != nil {
			logrus.Printf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			logrus.Printf("Server error: %v", err)
		}
	}
}

func handleGetColorInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hexCode, err := request.RequireString("hexCode")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	r, g, b, err := color.HexToRGB(hexCode)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	colorName := color.ClosestColorName(r, g, b, data.CssColors)
	result := map[string]string{
		"hex":  hexCode,
		"name": colorName,
		"rgb":  fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}
