package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	// Parse the command line flag "list"
	listFlag := flag.Bool("list", false, "list journal entries")
	flag.Parse()

	// If the "--list" flag was set, list the journal entries and exit
	if *listFlag {
		listEntries()
		return
	}

	// If the required number of arguments is not provided, display usage and exit
	if len(flag.Args()) < 3 {
		fmt.Println("Usage: notes <tags> <title> <description>")
		os.Exit(1)
	}

	// Parse the tags, title, and description from the command line arguments
	tags := strings.Split(flag.Arg(0), ",")
	title := flag.Arg(1)
	description := flag.Arg(2)

	// Create the journal directory if it doesn't exist
	journalDir := filepath.Join(os.Getenv("HOME"), "notes", "journal")
	if _, err := os.Stat(journalDir); os.IsNotExist(err) {
		os.MkdirAll(journalDir, 0755)
	}

	// Create the tags directory if it doesn't exist
	tagsDir := filepath.Join(journalDir, "tags")
	if _, err := os.Stat(tagsDir); os.IsNotExist(err) {
		os.MkdirAll(tagsDir, 0755)
	}

	// Create the new journal entry
	filename := title + ".txt"
	filedir := filepath.Join(journalDir, filename)
	file, err := os.Create(filedir)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write the description to the journal entry
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, description)
	w.Flush()

	// Create symlinks to the journal entry in the tags directories
	for _, tag := range tags {
		tagDir := filepath.Join(tagsDir, tag)
		if _, err := os.Stat(tagDir); os.IsNotExist(err) {
			os.MkdirAll(tagDir, 0755)
		}
		symlinkPath := filepath.Join(tagDir, filename)
		os.Symlink(filename, symlinkPath)
	}

	// Get the environment variable for the default file editor
	editor := os.Getenv("EDITOR")

	// If the environment variable is not set, use "nano" as the default editor
	if editor == "" {
		editor = "nano"
	}

	// Create a new command to open the specified file in the editor
	cmd := exec.Command(editor, filedir)

	// Set the command's standard input, output, and error to the current process's standard input, output, and error
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	cmd.Run()
}

func listEntries() {
	// Get the path to the journal directory
	journalDir := filepath.Join(os.Getenv("HOME"), "notes", "journal")

	// Get a list of all journal entries in the journal directory
	entries, err := filepath.Glob(filepath.Join(journalDir, "*.txt"))
	if err != nil {
		fmt.Printf("Error listing entries: %s\n", err)
		os.Exit(1)
	}

	// Iterate over the journal entries
	for i, entry := range entries {

		// Get the modification time and title of the entry
		info, err := os.Stat(entry)
		if err != nil {
			fmt.Printf("Error getting entry info: %s\n", err)
			os.Exit(1)
		}
		// Prepare the timestamp
		modTime := info.ModTime().Format("2006-01-02")
		title := strings.TrimSuffix(filepath.Base(entry), filepath.Ext(entry))

		// Get the size of the entry
		size, err := getFileSize(entry)
		if err != nil {
			fmt.Printf("Error getting entry size: %s\n", err)
			os.Exit(1)
		}

		// Display the entry
		fmt.Printf("(%d) %s %s %s\n", i+1, modTime, title, size)
	}

	// Prompt the user to enter the number of the entry they want to view
	fmt.Print("Enter the number of the entry you want to view: ")

	// Create a reader to read input from the standard input stream
	reader := bufio.NewReader(os.Stdin)

	// Read a line of input from the user
	input, _ := reader.ReadString('\n')

	// Trim leading and trailing whitespace from the input
	input = strings.TrimSpace(input)

	// Convert the input to an integer
	entryIndex, err := strconv.Atoi(input)

	// Check for an error in the conversion
	if err != nil {
		fmt.Printf("Error parsing input: %s\n", err)
		os.Exit(1)
	}

	// Check that the entry number is valid
	if entryIndex < 1 || entryIndex > len(entries) {
		fmt.Printf("Invalid entry number: %d\n", entryIndex)
		os.Exit(1)
	}

	// Open the selected journal entry in the default file editor
	entry := entries[entryIndex-1]
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}
	cmd := exec.Command(editor, entry)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

}

// Get the filesize in a human readable format
func getFileSize(filename string) (string, error) {
	// Get the file information
	info, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	// Get the size of the file in bytes
	size := info.Size()
	// If the size is less than 1024 bytes, return the size in bytes
	if size < 1024 {
		return fmt.Sprintf("%dB", size), nil
	} else if size < 1024*1024 { // If the size is less than 1024 * 1024 bytes, return the size in kilobytes
		return fmt.Sprintf("%.1fKB", float64(size)/1024), nil
	} else if size < 1024*1024*1024 { // If the size is less than 1024 * 1024 * 1024 bytes, return the size in megabytes
		return fmt.Sprintf("%.1fMB", float64(size)/1024/1024), nil
	}
	// If the size is more than or equal to 1024 * 1024 * 1024 bytes, return the size in gigabytes
	return fmt.Sprintf("%.1fGB", float64(size)/1024/1024/1024), nil
}
