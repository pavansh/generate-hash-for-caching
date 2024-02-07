package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Define a struct to hold the filename and its replacement pattern
type FileReplacement struct {
	Filename        string
	ReplacementHash string
}

func main() {
	var inputFile string
	var outputFile string
	var filesToReplace string
	var workDir string

	flag.StringVar(&inputFile, "input", "index.html", "Input HTML file")
	flag.StringVar(&outputFile, "output", "updated_index.html", "Output HTML file")
	flag.StringVar(&filesToReplace, "replace", "app.js,styles.css", "Comma-separated list of files to replace")
	flag.StringVar(&workDir, "workdir", ".", "Working directory")
	flag.Parse()

	// Change the working directory
	if err := os.Chdir(workDir); err != nil {
		fmt.Println("Error changing working directory:", err)
		os.Exit(1)
	}

	// Define the files to replace and their replacement patterns
	replacements := parseFileReplacements(filesToReplace)

	// Update the input HTML file
	err := updateIndexHTML(inputFile, outputFile, replacements)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Index.html updated successfully.")
}

func updateIndexHTML(inputFile, outputFile string, replacements []FileReplacement) error {
	// Read the content of the input HTML file
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Generate a new random seed
	seed := time.Now().UnixNano()

	ReplacementHash := generateRandomHash(seed)
	fmt.Println("Generated Hash: ", ReplacementHash)
	// Replace each file reference in the content
	updatedContent := string(content)
	for _, repl := range replacements {
		updatedContent = replaceFileReference(updatedContent, repl.Filename, ReplacementHash)
	}

	// Write the updated content to the output HTML file
	err = ioutil.WriteFile(outputFile, []byte(updatedContent), 0644)
	if err != nil {
		return err
	}

	// Rename each file in the replacements
	for _, repl := range replacements {
		err = renameFile(repl.Filename, ReplacementHash)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceFileReference(input, filename, hash string) string {
	// Define a regular expression pattern to match the file reference
	pattern := regexp.MustCompile(filename)
	// Replace the filename with the filename containing the hash
	return pattern.ReplaceAllStringFunc(input, func(match string) string {
		ext := filepath.Ext(match)                     // Get the extension of the matched filename
		base := match[:len(match)-len(ext)]            // Get the base name without the extension
		return fmt.Sprintf("%s.%s%s", base, hash, ext) // Replace with hash and original extension
	})
}

func generateRandomHash(seed int64) string {
	// Generate a random hash using SHA-256
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d", seed)))          // Use the provided seed
	randomHash := hex.EncodeToString(hash.Sum(nil))[:15] // Take the first 15 characters
	// Replace invalid characters (":", "/", "\" etc.) with "_"
	randomHash = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, randomHash)
	return randomHash
}

func renameFile(filename, hash string) error {
	// Split the filename into directory, base, and extension
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(filename)
	baseWithoutExt := base[:len(base)-len(ext)]
	// Construct the new filename with the hash
	newBase := fmt.Sprintf("%s.%s%s", baseWithoutExt, hash, ext)

	// Construct the full paths for the original and new files
	originalPath := filepath.Join(dir, base)
	newPath := filepath.Join(dir, newBase)

	// Rename the file
	if err := os.Rename(originalPath, newPath); err != nil {
		return fmt.Errorf("error renaming file %s to %s: %v", originalPath, newPath, err)
	}
	return nil
}

func parseFileReplacements(filesToReplace string) []FileReplacement {
	var replacements []FileReplacement
	files := strings.Split(filesToReplace, ",")
	for _, file := range files {
		replacements = append(replacements, FileReplacement{Filename: file})
	}
	return replacements
}
