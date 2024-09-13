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

const (
	journalDirName = "notes"
	tagsDirName    = "tags"
	defaultEditor  = "nano"
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

		// Display syntax and list entries
		fmt.Println("Usage: notes <tags> <title> <description>")
		listEntries()
		return
	}

	// Parse the tags, title, and description from the command line arguments
	tags := strings.Split(flag.Arg(0), ",")
	title := flag.Arg(1)
	description := flag.Arg(2)

	// Create the journal directory if it doesn't exist
	journalDir := getJournalDir()
	tagsDir := filepath.Join(journalDir, tagsDirName)
	createDirIfNotExist(journalDir)
	createDirIfNotExist(tagsDir)

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
		createDirIfNotExist(tagDir)
		symlinkPath := filepath.Join(tagDir, filename)
		if err := os.Symlink(filename, symlinkPath); err != nil {
			fmt.Printf("Error creating symlink: %s\n", err)
		}
	}

	// Get the environment variable for the default file editor
	editor := os.Getenv("EDITOR")

	// If the environment variable is not set, use "nano" as the default editor
	if editor == "" {
		editor = defaultEditor
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
	journalDir := getJournalDir()
	entries, err := filepath.Glob(filepath.Join(journalDir, "*.txt"))
	if err != nil {
		fmt.Printf("Error listing entries: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("\n<#> <date>\t<size>\t<title>")

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
		fmt.Printf("(%d) %s\t%s\t\"%s\"\n", i+1, modTime, size, title)
	}

	// Remove the redundant prompt and use promptForEntryNumber directly
	entryIndex := promptForEntryNumber(len(entries))

	// Open the selected journal entry in the default file editor
	entry := entries[entryIndex-1]
	openInEditor(entry)
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

func getJournalDir() string {
	return filepath.Join(os.Getenv("HOME"), journalDirName, "journal")
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory: %s\n", err)
			os.Exit(1)
		}
	}
}

func openInEditor(filepath string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error opening editor: %s\n", err)
		os.Exit(1)
	}
}

func promptForEntryNumber(maxEntries int) int {
	fmt.Print("\nEnter the number of the entry you want to view: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	entryIndex, err := strconv.Atoi(input)
	if err != nil || entryIndex < 1 || entryIndex > maxEntries {
		fmt.Printf("Invalid entry number: %s\n", input)
		os.Exit(1)
	}

	return entryIndex
}
