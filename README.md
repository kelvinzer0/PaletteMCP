# PaletteMCP - Color Code to Name Converter

PaletteMCP is a command-line tool written in Go that takes a hexadecimal color code and returns the name of the closest matching color from a predefined list of CSS colors. The output is provided in JSON format, making it easy to integrate with other scripts and systems.

## Features

- Convert any hex color code to its nearest color name.
- Outputs in a clean, machine-readable JSON format.
- Includes a comprehensive list of standard CSS colors.
- Simple and easy-to-use command-line interface.

## Installation

To use this tool, you need to have Go installed on your system.

1.  **Clone the repository (if you have it in a git repo) or just use the existing files.**

2.  **Build the executable:**
    ```sh
    go build -o palette-mcp ./cmd/palette-mcp
    ```

## Download Pre-built Binaries

You can download pre-built binaries for various operating systems and architectures directly from the [GitHub Releases page](https://github.com/kelvinzer0/PaletteMCP/releases).

Replace `[VERSION]` with the desired release version (e.g., `v1.0.0`).

### Linux / macOS

```bash
# Download the binary (replace [OS] and [ARCH] with your system, e.g., linux_amd64, darwin_arm64)
wget https://github.com/kelvinzer0/PaletteMCP/releases/download/[VERSION]/palette-mcp_[OS]_[ARCH] -O palette-mcp

# Make it executable
chmod +x palette-mcp

# Move it to a directory in your PATH (e.g., /usr/local/bin)
sudo mv palette-mcp /usr/local/bin/
```

### Windows

1.  Download the appropriate `.exe` file from the [GitHub Releases page](https://github.com/kelvinzer0/PaletteMCP/releases) (e.g., `palette-mcp_windows_amd64.exe`).
2.  Rename the downloaded file to `palette-mcp.exe`.
3.  Move `palette-mcp.exe` to a directory that is included in your system's `PATH` environment variable. A common practice is to create a `bin` folder in your user directory (e.g., `C:\Users\YourUser\bin`) and add it to `PATH`.



## Usage

`palette-mcp` can be used as a command-line tool or run as a Gemini Model Context Protocol (MCP) server.

### Command-Line Tool Usage

Run the tool from your terminal, passing a hex color code (with or without the `#` prefix) as an argument.

```sh
./palette-mcp #ff6347
```

### Example (Command-Line Tool)

**Input:**
```sh
./palette-mcp #ff6347
```

**Output:**
```json
{
  "hex": "#ff6347",
  "name": "Tomato",
  "rgb": "rgb(255, 99, 71)"
}
```

## Gemini MCP Server Integration

`palette-mcp` can also run as an MCP server, exposing its functionality to the Gemini CLI. This allows the Gemini model to discover and execute `palette-mcp`'s tools.

### Running the MCP Server

To start `palette-mcp` in server mode, use the `server` argument.

#### Stdio Transport

This is the default transport method and is used when you run the server with the `server` argument.

```sh
# Start the server using stdio
./palette-mcp server
```

#### HTTP Transport

You can also run the server with an HTTP transport, which allows you to specify a custom port.

```sh
# Start on default port 8080
./palette-mcp serve-http

# Start on a custom port (e.g., 9000)
./palette-mcp serve-http 9000
```

#### SSE (Server-Sent Events) Transport

To run the server in SSE mode, use the `-sse` flag. This will start an HTTP server that streams events.

```sh
# Start on default port 8080 in SSE mode
./palette-mcp -sse

# Start on a custom port (e.g., 9000) in SSE mode
./palette-mcp -sse -port 9000
```

It's recommended to run the server in the background if you want to continue using your terminal:

```sh
./palette-mcp serve-http &
# Or for a custom port:
./palette-mcp serve-http 9000 &

# For SSE mode:
./palette-mcp -sse &
# Or for a custom port:
./palette-mcp -sse -port 9000 & 
```

### Configuring Gemini CLI

To enable the Gemini CLI to connect to your `palette-mcp` server, add the following configuration to your `settings.json` file. This file can be found globally at `~/.gemini/settings.json` or in your project's `.gemini/settings.json`.

```json
{
  "theme": "ANSI Light",
  "selectedAuthType": "oauth-personal",
  "mcpServers": {
    "get-color-info": {
        "command": "/usr/local/bin/palette-mcp",
        "args":["server"]
    }
  }
}
```

*   **`paletteMcpServer`**: This is the name you give to your MCP server within Gemini CLI. You can choose any descriptive name.
*   **`httpUrl`**: The URL where your `palette-mcp` server is listening. Adjust the port if you started the server on a custom port (e.g., `http://localhost:9000`).
*   **`timeout`**: The maximum time (in milliseconds) Gemini CLI will wait for a response from the server.

### Available Tools

Once configured and the server is running, the Gemini CLI will discover the following tools:

*   **`echo`**: Echoes back the provided message. (Example tool for demonstration)
    *   **Parameters**: `message` (string)
*   **`get_color_info`**: Retrieves information about a color given its hex code.
    *   **Parameters**: `hexCode` (string, e.g., `#FF0000`)

### Example Gemini CLI Usage (Conceptual)

After setting up, you can interact with the tools via the Gemini CLI. For instance, you might ask:

`What is the name of the color #00BFFF?`

The Gemini model, if it decides to use the `get_color_info` tool, would execute it and provide the result.

## Integration with Forge MCP

`PaletteMCP` can be easily integrated into `forge mcp` workflows to provide color name lookup capabilities. You can add `palette-mcp` as a custom command within your `forge mcp` configuration.

### Adding PaletteMCP to Forge MCP

You can add `palette-mcp` to your `forge mcp` configuration using the `forge mcp add` command or by manually editing your `.mcp.json` file.

#### Using `forge mcp add` (Command Line)

```shell
# Example: Add palette-mcp as a command named 'colorname'
forge mcp add --name colorname --command /path/to/palette-mcp --args "#{{hex_code}}"
```
*   Replace `/path/to/palette-mcp` with the actual absolute path to your `palette-mcp` executable.
*   `#{{hex_code}}` is a placeholder for the hexadecimal color code that `forge mcp` will pass to `palette-mcp`. The double curly braces `{{...}}` indicate a variable or expression in `forge mcp`.

#### Manual `.mcp.json` Configuration

You can also manually create or modify your `.mcp.json` file. This file can be located in your local project directory or in your user-specific configuration.

```json
{
  "mcpServers": {
    "colorname_tool": {
      "command": "/path/to/palette-mcp",
      "args": ["#{{hex_code}}"],
      "description": "Converts a hex color code to its closest named color."
    }
  }
}
```
*   **`colorname_tool`**: This is the name you will use to invoke `palette-mcp` via `forge mcp` (e.g., `forge mcp run colorname_tool #RRGGBB`).
*   **`command`**: The absolute path to your `palette-mcp` executable.
*   **`args`**: An array of arguments to pass to `palette-mcp`. `"#{{hex_code}}"` is a common pattern for passing dynamic input from `forge mcp`.

### Example Forge MCP Usage

Once configured, you can use `palette-mcp` within your `forge mcp` workflows:

```shell
# Run the configured colorname_tool with a hex code
forge mcp run colorname_tool "#00BFFF"

# Example of using the output in a multi-agent workflow (conceptual)
# Assuming 'forge mcp run colorname_tool' outputs JSON
COLOR_INFO=$(forge mcp run colorname_tool "#FF0000")
COLOR_NAME=$(echo $COLOR_INFO | jq -r '.name')
echo "The color is: $COLOR_NAME"
```
*   **Note**: The `jq` command is a powerful JSON processor for the command line. You might need to install it separately (`brew install jq` on macOS, `sudo apt-get install jq` on Debian/Ubuntu).

This integration allows `forge mcp` to leverage `PaletteMCP` for color code to name conversions as part of larger automation or multi-agent tasks.



## How It Works

The tool calculates the "closest" color by finding the Euclidean distance between the input color's RGB values and the RGB values of each color in the predefined list. The color with the smallest distance is considered the closest match.

## Color Reference

The following table lists all the named colors used for matching.

| Color Name           | Hex     | RGB           |
| -------------------- | ------- | ------------- |
| aliceblue            | #f0f8ff | 240, 248, 255 |
| antiquewhite         | #faebd7 | 250, 235, 215 |
| aqua                 | #00ffff | 0, 255, 255   |
| aquamarine           | #7fffd4 | 127, 255, 212 |
| azure                | #f0ffff | 240, 255, 255 |
| beige                | #f5f5dc | 245, 245, 220 |
| bisque               | #ffe4c4 | 255, 228, 196 |
| black                | #000000 | 0, 0, 0       |
| blanchedalmond       | #ffebcd | 255, 235, 205 |
| blue                 | #0000ff | 0, 0, 255     |
| blueviolet           | #8a2be2 | 138, 43, 226  |
| brown                | #a52a2a | 165, 42, 42   |
| burlywood            | #deb887 | 222, 184, 135 |
| cadetblue            | #5f9ea0 | 95, 158, 160  |
| chartreuse           | #7fff00 | 127, 255, 0   |
| chocolate            | #d2691e | 210, 105, 30  |
| coral                | #ff7f50 | 255, 127, 80  |
| cornflowerblue       | #6495ed | 100, 149, 237 |
| cornsilk             | #fff8dc | 255, 248, 220 |
| crimson              | #dc143c | 220, 20, 60   |
| cyan                 | #00ffff | 0, 255, 255   |
| darkblue             | #00008b | 0, 0, 139     |
| darkcyan             | #008b8b | 0, 139, 139   |
| darkgoldenrod        | #b8860b | 184, 134, 11  |
| darkgray             | #a9a9a9 | 169, 169, 169 |
| darkgreen            | #006400 | 0, 100, 0     |
| darkkhaki            | #bdb76b | 189, 183, 107 |
| darkmagenta          | #8b008b | 139, 0, 139   |
| darkolivegreen       | #556b2f | 85, 107, 47   |
| darkorange           | #ff8c00 | 255, 140, 0   |
| darkorchid           | #9932cc | 153, 50, 204  |
| darkred              | #8b0000 | 139, 0, 0     |
| darksalmon           | #e9967a | 233, 150, 122 |
| darkseagreen         | #8fbc8f | 143, 188, 143 |
| darkslateblue        | #483d8b | 72, 61, 139   |
| darkslategray        | #2f4f4f | 47, 79, 79    |
| darkturquoise        | #00ced1 | 0, 206, 209   |
| darkviolet           | #9400d3 | 148, 0, 211   |
| deeppink             | #ff1493 | 255, 20, 147  |
| deepskyblue          | #00bfff | 0, 191, 255   |
| dimgray              | #696969 | 105, 105, 105 |
| dodgerblue           | #1e90ff | 30, 144, 255  |
| firebrick            | #b22222 | 178, 34, 34   |
| floralwhite          | #fffaf0 | 255, 250, 240 |
| forestgreen          | #228b22 | 34, 139, 34   |
| fuchsia              | #ff00ff | 255, 0, 255   |
| gainsboro            | #dcdcdc | 220, 220, 220 |
| ghostwhite              | #f8f8ff | 248, 248, 255 |
| gold                 | #ffd700 | 255, 215, 0   |
| goldenrod            | #daa520 | 218, 165, 32  |
| gray                 | #808080 | 128, 128, 128 |
| green                | #008000 | 0, 128, 0     |
| greenyellow          | #adff2f | 173, 255, 47  |
| honeydew             | #f0fff0 | 240, 255, 240 |
| hotpink              | #ff69b4 | 255, 105, 180 |
| indianred            | #cd5c5c | 205, 92, 92   |
| indigo               | #4b0082 | 75, 0, 130    |
| ivory                | #fffff0 | 255, 255, 240 |
| khaki                | #f0e68c | 240, 230, 140 |
| lavender             | #e6e6fa | 230, 230, 250 |
| lavenderblush        | #fff0f5 | 255, 240, 245 |
| lawngreen            | #7cfc00 | 124, 252, 0   |
| lemonchiffon         | #fffacd | 255, 250, 205 |
| lightblue            | #add8e6 | 173, 216, 230 |
| lightcoral           | #f08080 | 240, 128, 128 |
| lightcyan            | #e0ffff | 224, 255, 255 |
| lightgoldenrodyellow | #fafad2 | 250, 250, 210 |
| lightgray            | #d3d3d3 | 211, 211, 211 |
| lightgreen           | #90ee90 | 144, 238, 144 |
| lightpink            | #ffb6c1 | 255, 182, 193 |
| lightsalmon          | #ffa07a | 255, 160, 122 |
| lightseagreen        | #20b2aa | 32, 178, 170  |
| lightskyblue         | #87cefa | 135, 206, 250 |
| lightslategray       | #778899 | 119, 136, 153 |
| lightsteelblue       | #b0c4de | 176, 196, 222 |
| lightyellow          | #ffffe0 | 255, 255, 224 |
| lime                 | #00ff00 | 0, 255, 0     |
| limegreen            | #32cd32 | 50, 205, 50   |
| linen                | #faf0e6 | 250, 240, 230 |
| magenta              | #ff00ff | 255, 0, 255   |
| maroon               | #800000 | 128, 0, 0     |
| mediumaquamarine     | #66cdaa | 102, 205, 170 |
| mediumblue           | #0000cd | 0, 0, 205     |
| mediumorchid         | #ba55d3 | 186, 85, 211  |
| mediumpurple         | #9370db | 147, 112, 219 |
| mediumseagreen       | #3cb371 | 60, 179, 113  |
| mediumslateblue      | #7b68ee | 123, 104, 238 |
| mediumspringgreen    | #00fa9a | 0, 250, 154   |
| mediumturquoise      | #48d1cc | 72, 209, 204  |
| mediumvioletred      | #c71585 | 199, 21, 133  |
| midnightblue         | #191970 | 25, 25, 112   |
| mintcream            | #f5fffa | 245, 255, 250 |
| mistyrose            | #ffe4e1 | 255, 228, 225 |
| moccasin             | #ffe4b5 | 255, 228, 181 |
| navajowhite          | #ffdead | 255, 222, 173 |
| navy                 | #000080 | 0, 0, 128     |
| oldlace              | #fdf5e6 | 253, 245, 230 |
| olive                | #808000 | 128, 128, 0   |
| olivedrab            | #6b8e23 | 107, 142, 35  |
| orange               | #ffa500 | 255, 165, 0   |
| orangered            | #ff4500 | 255, 69, 0    |
| orchid               | #da70d6 | 218, 112, 214 |
| palegoldenrod        | #eee8aa | 238, 232, 170 |
| palegreen            | #98fb98 | 152, 251, 152 |
| paleturquoise        | #afeeee | 175, 238, 238 |
| palevioletred        | #db7093 | 219, 112, 147 |
| papayawhip           | #ffefd5 | 255, 239, 213 |
| peachpuff            | #ffdab9 | 255, 218, 185 |
| peru                 | #cd853f | 205, 133, 63  |
| pink                 | #ffc0cb | 255, 192, 203 |
| plum                 | #dda0dd | 221, 160, 221 |
| powderblue           | #b0e0e6 | 176, 224, 230 |
| purple               | #800080 | 128, 0, 128   |
| rebeccapurple        | #663399 | 102, 51, 153  |
| red                  | #ff0000 | 255, 0, 0     |
| rosybrown            | #bc8f8f | 188, 143, 143 |
| royalblue            | #4169e1 | 65, 105, 225  |
| saddlebrown          | #8b4513 | 139, 69, 19   |
| salmon               | #fa8072 | 250, 128, 114 |
| sandybrown           | #f4a460 | 244, 164, 96  |
| seagreen             | #2e8b57 | 46, 139, 87   |
| seashell             | #fff5ee | 255, 245, 238 |
| sienna               | #a0522d | 160, 82, 45   |
| silver               | #c0c0c0 | 192, 192, 192 |
| skyblue              | #87ceeb | 135, 206, 235 |
| slateblue            | #6a5acd | 106, 90, 205  |
| slategray            | #708090 | 112, 128, 144 |
| snow                 | #fffafa | 255, 250, 250 |
| springgreen          | #00ff7f | 0, 255, 127   |
| steelblue            | #4682b4 | 70, 130, 180  |
| tan                  | #d2b48c | 210, 180, 140 |
| teal                 | #008080 | 0, 128, 128   |
| thistle              | #d8bfd8 | 216, 191, 216 |
| tomato               | #ff6347 | 255, 99, 71   |
| turquoise            | #40e0d0 | 64, 224, 208  |
| violet               | #ee82ee | 238, 130, 238 |
| wheat                | #f5deb3 | 245, 222, 179 |
| white                | #ffffff | 255, 255, 255 |
| whitesmoke           | #f5f5f5 | 245, 245, 245 |
| yellow               | #ffff00 | 255, 255, 0   |
| yellowgreen          | #9acd32 | 154, 205, 50  |

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have suggestions for improvements.

## License

This project is licensed under the MIT License.