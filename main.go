package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	journalDirName = "notes"
	tagsDirName    = "tags"
	defaultEditor  = "nano"
)

func main() {
	// Create command
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createTitle := createCmd.String("title", "", "title of the note (required)")
	createBody := createCmd.String("body", "", "body of the note")
	createTags := createCmd.String("tags", "", "comma-separated list of tags")

	// List command
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listTag := listCmd.String("tag", "", "list entries for a specific tag")

	// Tags command
	tagsCmd := flag.NewFlagSet("tags", flag.ExitOnError)
	tagsSearch := tagsCmd.String("search", "", "search entries by tag")

	if len(os.Args) == 2 && os.Args[1] == "--help" {
		printUsage()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("Expected 'create', 'list', or 'tags' subcommands")
		fmt.Println("Use --help for more information")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if createCmd.Parsed() {
			if *createTitle == "" && len(createCmd.Args()) == 0 {
				fmt.Println("Usage: journal create [options]")
				createCmd.PrintDefaults()
				os.Exit(0)
			}
			createNote(*createTitle, *createBody, *createTags)
		}
	case "list":
		listCmd.Parse(os.Args[2:])
		if listCmd.Parsed() {
			if len(listCmd.Args()) == 0 && *listTag == "" {
				fmt.Println("Usage: journal list [options]")
				listCmd.PrintDefaults()
				os.Exit(0)
			}
			if *listTag != "" {
				listEntriesByTag(*listTag)
			} else {
				listEntries()
			}
		}
	case "tags":
		tagsCmd.Parse(os.Args[2:])
		if tagsCmd.Parsed() {
			if len(tagsCmd.Args()) == 0 && *tagsSearch == "" {
				fmt.Println("Usage: journal tags [options]")
				tagsCmd.PrintDefaults()
				os.Exit(0)
			}
			if *tagsSearch != "" {
				searchEntriesByTag(*tagsSearch)
			} else {
				listAllTags()
			}
		}
	default:
		fmt.Println("Expected 'create', 'list', or 'tags' subcommands")
		fmt.Println("Use --help for more information")
		os.Exit(1)
	}
}

func createNote(title, body, tags string) {
	if title == "" {
		fmt.Println("Error: Title is required")
		os.Exit(1)
	}

	journalDir := getJournalDir()
	tagsDir := filepath.Join(journalDir, tagsDirName)
	createDirIfNotExist(journalDir)
	createDirIfNotExist(tagsDir)

	filename := title + ".txt"
	filedir := filepath.Join(journalDir, filename)
	file, err := os.Create(filedir)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, body)
	w.Flush()

	tagList := strings.Split(tags, ",")
	for _, tag := range tagList {
		tag = strings.TrimSpace(tag)
		tagDir := filepath.Join(tagsDir, tag)
		createDirIfNotExist(tagDir)
		symlinkPath := filepath.Join(tagDir, filename)
		if err := os.Symlink(filename, symlinkPath); err != nil {
			fmt.Printf("Error creating symlink: %s\n", err)
		}
	}

	openInEditor(filedir)
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

func listAllTags() {
	tagsDir := filepath.Join(getJournalDir(), tagsDirName)
	entries, err := os.ReadDir(tagsDir)
	if err != nil {
		fmt.Printf("Error listing tags: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("All tags:")
	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip non-directory entries
		}
		name := entry.Name()
		if !isHiddenOrSystemFile(name) {
			fmt.Println(name)
		}
	}
}

func listEntriesByTag(tag string) {
	tagDir := filepath.Join(getJournalDir(), tagsDirName, tag)
	entries, err := filepath.Glob(filepath.Join(tagDir, "*.txt"))
	if err != nil {
		fmt.Printf("Error listing entries for tag '%s': %s\n", tag, err)
		os.Exit(1)
	}

	fmt.Printf("Entries tagged with '%s':\n", tag)
	displayEntries(entries)
}

func searchEntriesByTag(searchTag string) {
	tagsDir := filepath.Join(getJournalDir(), tagsDirName)
	entries, err := os.ReadDir(tagsDir)
	if err != nil {
		fmt.Printf("Error searching tags: %s\n", err)
		os.Exit(1)
	}

	var matchingEntries []string
	for _, entry := range entries {
		if !entry.IsDir() || isHiddenOrSystemFile(entry.Name()) {
			continue
		}
		if strings.Contains(entry.Name(), searchTag) {
			tagDir := filepath.Join(tagsDir, entry.Name())
			files, err := filepath.Glob(filepath.Join(tagDir, "*.txt"))
			if err == nil {
				matchingEntries = append(matchingEntries, files...)
			}
		}
	}

	fmt.Printf("Entries matching tag search '%s':\n", searchTag)
	displayEntries(matchingEntries)
}

func displayEntries(entries []string) {
	sort.Strings(entries)
	for i, entry := range entries {
		info, err := os.Stat(entry)
		if err != nil {
			continue
		}
		modTime := info.ModTime().Format("2006-01-02")
		title := strings.TrimSuffix(filepath.Base(entry), filepath.Ext(entry))
		size, _ := getFileSize(entry)
		fmt.Printf("(%d) %s\t%s\t\"%s\"\n", i+1, modTime, size, title)
	}

	if len(entries) > 0 {
		entryIndex := promptForEntryNumber(len(entries))
		openInEditor(entries[entryIndex-1])
	} else {
		fmt.Println("No entries found.")
	}
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

func isHiddenOrSystemFile(name string) bool {
	// List of common hidden or system directories to ignore
	ignoredPrefixes := []string{".", "_", "~$"}
	ignoredNames := map[string]bool{
		"System Volume Information": true,
		"$RECYCLE.BIN":              true,
		"lost+found":                true,
	}

	// Check if the name starts with any ignored prefix
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	// Check if the name is in the list of ignored names
	return ignoredNames[name]
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  journal [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  create  Create a new journal entry")
	fmt.Println("  list    List journal entries")
	fmt.Println("  tags    List or search tags")
	fmt.Println("\nOptions:")
	fmt.Println("  --help  Show help information for a command")
	fmt.Println("\nUse 'journal [command] --help' for more information about a command.")
}
