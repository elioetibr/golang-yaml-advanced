// Package internal provides internal helper functions for the yaml package.
// These functions are not part of the public API and may change without notice.
//
// This package contains only generic utilities that don't depend on the Node type.
// Node-specific helpers remain as unexported functions in the yaml package.
package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidateFormat validates a string value against a specified format.
// Supported formats: email, uri, url, ipv4, ipv6, uuid, date, date-time, datetime, time.
func ValidateFormat(value, format string) bool {
	switch format {
	case "email":
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+$`
		matched, _ := regexp.MatchString(emailRegex, value)
		return matched
	case "uri", "url":
		urlRegex := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
		matched, _ := regexp.MatchString(urlRegex, value)
		return matched
	case "ipv4":
		ipv4Regex := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
		matched, _ := regexp.MatchString(ipv4Regex, value)
		return matched
	case "ipv6":
		ipv6Regex := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:))$`
		matched, _ := regexp.MatchString(ipv6Regex, value)
		return matched
	case "uuid":
		uuidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
		matched, _ := regexp.MatchString(uuidRegex, value)
		return matched
	case "date":
		dateRegex := `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`
		matched, _ := regexp.MatchString(dateRegex, value)
		if !matched {
			return false
		}
		// Additional date validation
		parts := strings.Split(value, "-")
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		day, _ := strconv.Atoi(parts[2])
		// Check days in month
		daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			daysInMonth[1] = 29 // Leap year
		}
		if month < 1 || month > 12 || day < 1 || day > daysInMonth[month-1] {
			return false
		}
		return true
	case "date-time", "datetime":
		dateTimeRegex := `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](\.\d+)?(Z|[+-]\d{2}:\d{2})$`
		matched, _ := regexp.MatchString(dateTimeRegex, value)
		return matched
	case "time":
		timeRegex := `^([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`
		matched, _ := regexp.MatchString(timeRegex, value)
		return matched
	default:
		return true // Unknown format, assume valid
	}
}

// SplitDocuments splits a YAML string into separate documents by --- separator.
// This function handles multi-document YAML files.
func SplitDocuments(content string) []string {
	lines := strings.Split(content, "\n")
	var documents []string
	var currentDoc strings.Builder
	inDocument := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "---" {
			// Document separator found
			if currentDoc.Len() > 0 {
				documents = append(documents, currentDoc.String())
				currentDoc.Reset()
			}
			inDocument = true
		} else if trimmed == "..." {
			// Document end marker
			if currentDoc.Len() > 0 {
				documents = append(documents, currentDoc.String())
				currentDoc.Reset()
			}
			inDocument = false
		} else {
			// Regular content line
			if !inDocument && len(documents) == 0 {
				// First document without explicit --- marker
				inDocument = true
			}
			if inDocument {
				if currentDoc.Len() > 0 {
					currentDoc.WriteString("\n")
				}
				currentDoc.WriteString(line)
			}
		}
	}

	// Add the last document if any
	if currentDoc.Len() > 0 {
		documents = append(documents, currentDoc.String())
	}

	// If no documents were found, treat the entire content as one document
	if len(documents) == 0 && len(content) > 0 {
		documents = append(documents, content)
	}

	return documents
}

// InterfaceToMap converts various types to map[string]interface{}.
// This is used for flexible merging operations.
func InterfaceToMap(input interface{}) (map[string]interface{}, error) {
	// Already a map
	if m, ok := input.(map[string]interface{}); ok {
		return m, nil
	}
	if m, ok := input.(map[interface{}]interface{}); ok {
		result := make(map[string]interface{})
		for k, v := range m {
			if str, ok := k.(string); ok {
				result[str] = v
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("cannot convert %T to map[string]interface{}", input)
}

// MergeMaps recursively merges two maps.
// For conflicting keys, override values take precedence unless both are maps,
// in which case they are recursively merged.
func MergeMaps(base, override map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy base values
	for k, v := range base {
		result[k] = v
	}

	// Apply overrides
	for k, overrideValue := range override {
		baseValue, exists := result[k]

		if !exists {
			// Key doesn't exist in base, add it
			result[k] = overrideValue
			continue
		}

		// Both exist - try to merge if both are maps
		if baseMap, baseIsMap := baseValue.(map[string]interface{}); baseIsMap {
			if overrideMap, overrideIsMap := overrideValue.(map[string]interface{}); overrideIsMap {
				result[k] = MergeMaps(baseMap, overrideMap)
				continue
			}
		}

		// For arrays, append override to base
		if baseSlice, baseIsSlice := baseValue.([]interface{}); baseIsSlice {
			if overrideSlice, overrideIsSlice := overrideValue.([]interface{}); overrideIsSlice {
				result[k] = append(baseSlice, overrideSlice...)
				continue
			}
		}

		// Otherwise, override takes precedence
		result[k] = overrideValue
	}

	return result
}

// AddEmptyLinesBeforeCommentBlocks adds empty lines before comment blocks.
// This uses heuristics to preserve formatting conventions.
func AddEmptyLinesBeforeCommentBlocks(input []byte) []byte {
	lines := strings.Split(string(input), "\n")
	var output []string
	inCommentBlock := false

	for i := 0; i < len(lines); i++ {
		currentLine := lines[i]
		trimmedCurrent := strings.TrimSpace(currentLine)

		// Determine if current line is a comment
		isComment := strings.HasPrefix(trimmedCurrent, "#")

		// Check if we need to add an empty line before this line
		if i > 0 {
			prevLine := lines[i-1]
			trimmedPrev := strings.TrimSpace(prevLine)
			prevIsComment := strings.HasPrefix(trimmedPrev, "#")

			// Add empty line before a comment block that follows a non-comment line
			if isComment && !prevIsComment && trimmedPrev != "" {
				// This is the start of a new comment block after content
				output = append(output, "")
				inCommentBlock = true
			} else if isComment && prevIsComment {
				// Continue in comment block
				inCommentBlock = true
			} else if !isComment && inCommentBlock {
				// End of comment block
				inCommentBlock = false
			}
		} else if isComment {
			inCommentBlock = true
		}

		output = append(output, currentLine)
	}

	return []byte(strings.Join(output, "\n"))
}

// RemoveEmptyLines removes all empty lines from the output.
func RemoveEmptyLines(input []byte) []byte {
	lines := strings.Split(string(input), "\n")
	var output []string

	for _, line := range lines {
		// Keep non-empty lines
		if strings.TrimSpace(line) != "" {
			output = append(output, line)
		}
	}

	// Add a final newline if there are lines
	if len(output) > 0 {
		return []byte(strings.Join(output, "\n") + "\n")
	}
	return []byte{}
}

// NormalizeEmptyLines ensures consistent empty line counts.
// It normalizes empty lines between sections to the specified count.
func NormalizeEmptyLines(input []byte, count int) []byte {
	lines := strings.Split(string(input), "\n")
	var output []string
	lastWasContent := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		isComment := strings.HasPrefix(trimmed, "#")
		isEmpty := trimmed == ""

		if isEmpty {
			// Skip empty lines, we'll add them as needed
			continue
		}

		// Add normalized empty lines before comments (except at start)
		if isComment && lastWasContent && i > 0 {
			for j := 0; j < count; j++ {
				output = append(output, "")
			}
		}

		output = append(output, line)
		lastWasContent = !isComment
	}

	return []byte(strings.Join(output, "\n"))
}
