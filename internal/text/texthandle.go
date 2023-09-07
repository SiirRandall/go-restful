package text // Package 'text' handles visualizing and formatting text

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// VisualizeJSONStructure is a function to create a formatted string representation of JSON data
func VisualizeJSONStructure(data interface{}, indent string, jsonwidth int) string {
	switch v := data.(type) { // A type switch to handle different types of JSON structures
	case map[string]interface{}: // If the data is a JSON object
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "{\n")

		// Extract and sort the keys for consistent output.
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys) // Sorting the keys for consistent order

		for i, key := range keys {
			valueStr := VisualizeJSONStructure(v[key], indent+"  ", jsonwidth)
			comma := ","
			if i == len(keys)-1 {
				comma = "" // No comma for the last key
			}
			// Write each key-value pair with proper indentation, coloring, and trailing comma (if not last item)
			fmt.Fprintf(b, "%s  [blue]%s[white]: %s%s\n", indent, key, strings.Trim(valueStr, "\n"), comma)
		}

		fmt.Fprintf(b, "%s}", indent)
		return b.String()

	case []interface{}: // If the data is a JSON array
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "[\n")

		for i, item := range v {
			itemStr := VisualizeJSONStructure(item, indent+"  ", jsonwidth)
			comma := ","
			if i == len(v)-1 {
				comma = "" // No comma for the last item
			}
			// Write each item with proper indentation, coloring, and trailing comma (if not last item)
			fmt.Fprintf(b, "%s  %s%s\n", indent, strings.Trim(itemStr, "\n"), comma)
		}

		fmt.Fprintf(b, "%s]", indent)
		return b.String()

	default:
		// Apply custom styles based on the type of the value
		switch t := data.(type) {
		case string: // If the data is a string
			t = strings.ReplaceAll(t, "\n", "\n") // Replace newline characters with \n literal

			// Word wrap for strings exceeding the JSON view width
			t = WordWrap(t, jsonwidth-len(indent)-2) // Subtract 2 for quotes
			// Indent every new line after wrapping
			t = strings.ReplaceAll(t, "\n", "\n"+indent+"  ")

			return fmt.Sprintf("[orange]\"%s\"[white]", t)
		case float64: // If the data is a floating point number
			if t == float64(int(t)) { // Check if the number is an integer
				return fmt.Sprintf("[red]%d[white]", int(t))
			}
			return fmt.Sprintf("[green]%f[white]", t)
		case int, int32, int64: // If the data is an integer
			return fmt.Sprintf("[red]%d[white]", t)
		default: // For any other type
			return fmt.Sprintf("%v", t)
		}
	}
}

// WordWrap is a function to ensure that a given text string wraps at a specified width
func WordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	var buffer bytes.Buffer
	words := strings.Fields(text) // Split the input text into words

	if len(words) == 0 {
		return ""
	}

	buffer.WriteString(words[0]) // Write the first word to the buffer
	spaceLeft := width - len(words[0])

	for _, word := range words[1:] { // For each remaining word in the text
		if len(word)+1 > spaceLeft { // If writing this word would exceed the max line width...
			buffer.WriteString("\n" + word) // ...then write it on a new line
			spaceLeft = width - len(word)
		} else {
			buffer.WriteString(" " + word) // Otherwise, write it on the same line
			spaceLeft -= (1 + len(word))
		}
	}

	return buffer.String()
}
