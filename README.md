

# log-beautify

A command-line tool that makes logs more readable by formatting and prettifying them, with special handling for JSON and multi-line logs.

> I developed this tool using the Cursor, an incredibly powerful AI programming tool that helps me develop tools of all kinds, fast.

## Features

- **JSON Log Processing**
  - Pretty prints JSON with proper indentation
  - Automatically unescapes string values for better readability
  - Shows escaped values (like `\n`, `\"`) in a human-readable format

- **Multi-line Log Support**
  - Handles indentation-based multi-line logs
  - Treats indented lines as continuation of the previous log entry
  - Supports empty-line separated log entries

- **Flexible Output**
  - Terminal-friendly output with syntax highlighting
  - Optional raw markdown output for further processing
  - Clean and consistent formatting

## Installation

```bash
go install github.com/mylxsw/log-beautify@latest
```

## Usage

### Basic Usage

Process logs from a file:

```bash
cat logfile.log | log-beautify
```

Process logs from clipboard (macOS):

```bash
pbpaste | log-beautify
```

### Options

```bash
log-beautify -raw    # Output raw markdown without terminal rendering
```

## Examples

### JSON Log Processing

Input:
```json
{"message": "Hello\\nWorld\\\"Test\\\"", "level": "info"}
```

Output:
```
## JSON Log

{
  "message": "--- SEE BELOW ---",
  "level": "info"
}

### Field `message` Unescaped Value:

Hello
World"Test"
```

### Multi-line Log Processing

Input:
```
Starting process
  with additional details
  and more information
Another log entry
```

Output:
```
## Plain Text Log

Starting process
  with additional details
  and more information

## Plain Text Log

Another log entry
```

## How It Works

1. Reads input stream line by line
2. Groups related log lines based on indentation
3. Attempts to parse each log entry as JSON
4. For JSON logs:
   - Pretty prints the structure
   - Extracts and unescapes string values
5. For plain text:
   - Preserves original formatting
   - Groups indented lines with their parent
6. Formats output using markdown
7. Renders for terminal display (unless -raw is specified)

## Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests

## License

MIT
