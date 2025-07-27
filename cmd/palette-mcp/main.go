package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/kelvinzer0/palette-mcp/internal/color"
	"github.com/kelvinzer0/palette-mcp/internal/data"
)

func main() {
	// Parse command line flags
	sseMode := flag.Bool("sse", false, "Run in SSE mode instead of stdio mode")
	port := flag.String("port", "8080", "Port for SSE or HTTP server")
	flag.Parse()

	if *sseMode {
		startMcpServer(false, *port, true) // SSE mode
	} else if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			startMcpServer(false, "", false) // Stdio mode
		case "serve-http":
			// If serve-http is used, check for a port argument after it
			if len(os.Args) > 2 {
				*port = os.Args[2]
			}
			startMcpServer(true, *port, false) // HTTP mode
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

func startMcpServer(useHTTP bool, port string, useSSE bool) {
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

	if useSSE {
		sseServer := server.NewSSEServer(s, server.WithBaseURL("http://localhost:" + port))
		logrus.Printf("Starting SSE server on localhost:%s", port)
		if err := sseServer.Start(":" + port); err != nil {
			logrus.Fatalf("Server error: %v", err)
		}
	} else if useHTTP {
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
