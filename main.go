package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
)

// Process a single log entry and return markdown text
func processLog(log string) string {
	var output strings.Builder

	// Try to parse as JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(log), &jsonData); err == nil {
		output.WriteString("## JSON Log\n\n")
		// Process values for display before output
		processedJSON := processJSONForDisplay(jsonData)
		// Pretty print JSON
		prettyJSON, _ := json.MarshalIndent(processedJSON, "", "  ")
		// 处理多行文本，确保每行都正确缩进
		jsonLines := strings.Split(string(prettyJSON), "\n")
		output.WriteString("```json\n")
		for _, line := range jsonLines {
			output.WriteString(line)
			output.WriteString("\n")
		}
		output.WriteString("```\n")
		// Collect unescaped values
		collectEscapedValues(&output, jsonData)
	} else {
		output.WriteString("## Plain Text Log\n\n")
		output.WriteString("```\n")
		output.WriteString(log)
		output.WriteString("\n```\n")
	}
	output.WriteString("\n---\n\n")
	return output.String()
}

// Process JSON for display
func processJSONForDisplay(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			// 检查字符串是否包含换行符（包括转义的换行符）
			unescaped := strings.ReplaceAll(v, "\\\\", "\\")
			unescaped = strings.ReplaceAll(unescaped, "\\n", "\n")
			unescaped = strings.ReplaceAll(unescaped, "\\r", "\r")
			unescaped = strings.ReplaceAll(unescaped, "\\\"", "\"")
			if strings.Contains(unescaped, "\n") || len(unescaped) > 100 {
				result[key] = "--- SEE BELOW ---"
			} else {
				result[key] = v
			}
		case []interface{}:
			allStrings := true
			shouldFormat := false
			for _, item := range v {
				if str, ok := item.(string); ok {
					unescaped := strings.ReplaceAll(str, "\\\\", "\\")
					unescaped = strings.ReplaceAll(unescaped, "\\n", "\n")
					unescaped = strings.ReplaceAll(unescaped, "\\r", "\r")
					unescaped = strings.ReplaceAll(unescaped, "\\\"", "\"")
					if strings.Contains(unescaped, "\n") {
						shouldFormat = true
					}
				} else {
					allStrings = false
					break
				}
			}
			if allStrings && shouldFormat {
				result[key] = "--- SEE BELOW ---"
			} else {
				result[key] = v
			}
		case map[string]interface{}:
			result[key] = processJSONForDisplay(v)
		default:
			result[key] = v
		}
	}
	return result
}

// Collect escaped values from JSON
func collectEscapedValues(output *strings.Builder, data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			unescaped := strings.ReplaceAll(v, "\\\\", "\\")
			unescaped = strings.ReplaceAll(unescaped, "\\n", "\n")
			unescaped = strings.ReplaceAll(unescaped, "\\r", "\r")
			unescaped = strings.ReplaceAll(unescaped, "\\\"", "\"")
			// 如果字符串包含换行符或长度超过100，则单独显示
			if strings.Contains(unescaped, "\n") || len(unescaped) > 100 {
				fmt.Fprintf(output, "\n### Field `%s`:\n\n```\n%s\n```\n", key, unescaped)
			}
		case []interface{}:
			allStrings := true
			var stringsToFormat []string
			shouldFormat := false

			for _, item := range v {
				if str, ok := item.(string); ok {
					unescaped := strings.ReplaceAll(str, "\\\\", "\\")
					unescaped = strings.ReplaceAll(unescaped, "\\n", "\n")
					unescaped = strings.ReplaceAll(unescaped, "\\r", "\r")
					unescaped = strings.ReplaceAll(unescaped, "\\\"", "\"")
					if strings.Contains(unescaped, "\n") {
						shouldFormat = true
					}
					stringsToFormat = append(stringsToFormat, unescaped)
				} else {
					allStrings = false
					break
				}
			}

			if allStrings && shouldFormat {
				fmt.Fprintf(output, "\n### Field `%s` Values:\n\n", key)
				for i, str := range stringsToFormat {
					fmt.Fprintf(output, "#### Item %d:\n```\n%s\n```\n", i+1, str)
				}
			}
		case map[string]interface{}:
			collectEscapedValues(output, v)
		}
	}
}

func main() {
	// Add command line flag
	noRender := flag.Bool("raw", false, "Output raw markdown text without rendering")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	var currentLog strings.Builder
	var lastLineWasIndented bool
	var allOutput strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// If empty line, process collected log
		if line == "" {
			if currentLog.Len() > 0 {
				allOutput.WriteString(processLog(currentLog.String()))
				currentLog.Reset()
			}
			lastLineWasIndented = false
			continue
		}

		// Check if line starts with whitespace
		isIndented := len(line) > 0 && (line[0] == ' ' || line[0] == '\t')

		// If this is a new non-indented line and we have content, process previous log
		if !isIndented && currentLog.Len() > 0 && !lastLineWasIndented {
			allOutput.WriteString(processLog(currentLog.String()))
			currentLog.Reset()
		}

		// Add current line to log
		if currentLog.Len() > 0 {
			currentLog.WriteString("\n")
		}
		currentLog.WriteString(line)
		lastLineWasIndented = isIndented
	}

	// Process final log entry
	if currentLog.Len() > 0 {
		allOutput.WriteString(processLog(currentLog.String()))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return
	}

	// Choose output method based on flag
	if *noRender {
		fmt.Print(allOutput.String())
	} else {
		result := markdown.Render(allOutput.String(), 80, 0)
		fmt.Print(string(result))
	}
}
