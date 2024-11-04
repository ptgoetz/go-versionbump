package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// ReplaceInFile replaces all occurrences of the search string with the replace string in the file at the given path.
func ReplaceInFile(filePath string, search string, replace string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file line by line and replace the version string
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		updatedLine := strings.ReplaceAll(line, search, replace)
		lines = append(lines, updatedLine)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write the updated lines back to the file
	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceInString(input string, search string, replace string) string {
	return strings.ReplaceAll(input, search, replace)
}

// CountStringsInFile returns the number of times the search string occurs in the file at the given path.
func CountStringsInFile(filePath, searchString string) (int, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Initialize a counter for the occurrences
	count := 0

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Count the occurrences of the search string in the current line
		count += strings.Count(line, searchString)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

// ParentDirAbsolutePath returns the absolute path of the parent directory of the given relative file path.
func ParentDirAbsolutePath(relativePath string) (string, error) {
	// Get the absolute path of the file
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}

	// Get the parent directory
	parentDir := filepath.Dir(absolutePath)

	// Return the absolute path of the parent directory
	return parentDir, nil
}

// PaddingString returns a string of the specified length filled with the given padding character.
func PaddingString(length int, padChar string) string {
	if len(padChar) != 1 {
		panic("padChar must be a single character")
	}
	if length <= 0 {
		return ""
	}
	return strings.Repeat(padChar, length)
}

// IsAllAlphabetic returns true if the given string contains only alphabetic characters.
func IsAllAlphabetic(s string) bool {
	for _, char := range s {
		// Check if the character is not a letter or is non-ASCII
		if !unicode.IsLetter(char) || char > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsAllAlphanumeric returns true if the given string contains only alphanumeric characters.
func IsAllAlphanumeric(s string) bool {
	for _, char := range s {
		// Check if the character is not a letter, digit, or is non-ASCII
		if !(unicode.IsLetter(char) || unicode.IsDigit(char)) || char > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// StartsWithDigit returns true if the given string starts with a digit.
func StartsWithDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsDigit(rune(s[0]))
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
