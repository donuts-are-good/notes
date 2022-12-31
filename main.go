package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	listFlag := flag.Bool("list", false, "list journal entries")
	flag.Parse()

	if *listFlag {
		listEntries()
		return
	}

	if len(flag.Args()) < 3 {
		fmt.Println("Usage: notes <tags> <title> <description>")
		os.Exit(1)
	}

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

	/// Create symlinks to the journal entry in the tags directories
	for _, tag := range tags {
		tagDir := filepath.Join(tagsDir, tag)
		if _, err := os.Stat(tagDir); os.IsNotExist(err) {
			os.MkdirAll(tagDir, 0755)
		}
		symlinkPath := filepath.Join(tagDir, filename)
		os.Symlink(filename, symlinkPath)
	}

	// Open the journal entry in the default file editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, filedir)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func listEntries() {
	journalDir := filepath.Join(os.Getenv("HOME"), "notes", "journal")
	entries, err := filepath.Glob(filepath.Join(journalDir, "*.txt"))
	if err != nil {
		fmt.Printf("Error listing entries: %s\n", err)
		os.Exit(1)
	}

	for i, entry := range entries {
		fmt.Printf("%d) %s\n", i+1, filepath.Base(entry))
	}

	fmt.Print("Enter the number of the entry you want to view: ")
	var n int
	_, err = fmt.Scanf("%d", &n)
	if err != nil || n < 1 || n > len(entries) {
		fmt.Println("Invalid input")
		os.Exit(1)
	}

	// Open the selected journal entry in the default file editor
	entry := entries[n-1]
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, entry)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
