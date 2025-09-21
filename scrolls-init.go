package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Screenshot  string    `json:"screenshot,omitempty"`
}

type NoteManager struct {
	dataDir string
	notes   []Note
	nextID  int
}

func NewNoteManager() *NoteManager {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".ancient-scrolls")
	
	// Create data directory if it doesn't exist
	os.MkdirAll(dataDir, 0755)
	
	nm := &NoteManager{
		dataDir: dataDir,
		notes:   []Note{},
		nextID:  1,
	}
	
	nm.loadNotes()
	return nm
}

func (nm *NoteManager) loadNotes() {
	files, err := ioutil.ReadDir(nm.dataDir)
	if err != nil {
		return
	}
	
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			data, err := ioutil.ReadFile(filepath.Join(nm.dataDir, file.Name()))
			if err != nil {
				continue
			}
			
			var note Note
			if err := json.Unmarshal(data, &note); err != nil {
				continue
			}
			
			nm.notes = append(nm.notes, note)
			if note.ID >= nm.nextID {
				nm.nextID = note.ID + 1
			}
		}
	}
	
	// Sort notes by creation time
	sort.Slice(nm.notes, func(i, j int) bool {
		return nm.notes[i].CreatedAt.After(nm.notes[j].CreatedAt)
	})
}

func (nm *NoteManager) saveNote(note Note) error {
	filename := filepath.Join(nm.dataDir, fmt.Sprintf("note_%d.json", note.ID))
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(filename, data, 0644)
}

func (nm *NoteManager) createNote() {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("Enter note title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	
	if title == "" {
		fmt.Println("Title cannot be empty.")
		return
	}
	
	fmt.Print("Enter note content (press Ctrl+D when finished):\n")
	var content strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}
	
	fmt.Print("Enter tags (comma-separated): ")
	tagsInput, _ := reader.ReadString('\n')
	tagsInput = strings.TrimSpace(tagsInput)
	
	var tags []string
	if tagsInput != "" {
		for _, tag := range strings.Split(tagsInput, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}
	
	fmt.Print("Take screenshot? (y/n): ")
	screenshotChoice, _ := reader.ReadString('\n')
	screenshotChoice = strings.TrimSpace(screenshotChoice)
	
	var screenshot string
	if strings.ToLower(screenshotChoice) == "y" {
		screenshot = nm.takeScreenshot()
	}
	
	note := Note{
		ID:         nm.nextID,
		Title:      title,
		Content:    strings.TrimSpace(content.String()),
		Tags:       tags,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Screenshot: screenshot,
	}
	
	if err := nm.saveNote(note); err != nil {
		fmt.Printf("Error saving note: %v\n", err)
		return
	}
	
	nm.notes = append([]Note{note}, nm.notes...)
	nm.nextID++
	
	fmt.Printf("Note created successfully with ID: %d\n", note.ID)
}

func (nm *NoteManager) takeScreenshot() string {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("screenshot_%s.png", timestamp)
	filepath := filepath.Join(nm.dataDir, filename)
	
	// Try different screenshot commands based on what's available
	commands := [][]string{
		{"gnome-screenshot", "-f", filepath},
		{"scrot", filepath},
		{"import", "-window", "root", filepath},
	}
	
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd[0]); err == nil {
			if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
				fmt.Printf("Screenshot saved: %s\n", filename)
				return filename
			}
		}
	}
	
	fmt.Println("No screenshot tool found. Install gnome-screenshot, scrot, or imagemagick.")
	return ""
}

func (nm *NoteManager) listNotes() {
	if len(nm.notes) == 0 {
		fmt.Println("No notes found.")
		return
	}
	
	fmt.Printf("%-4s %-30s %-20s %-15s\n", "ID", "Title", "Created", "Tags")
	fmt.Println(strings.Repeat("-", 70))
	
	for _, note := range nm.notes {
		tags := strings.Join(note.Tags, ", ")
		if len(tags) > 15 {
			tags = tags[:12] + "..."
		}
		
		title := note.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		
		fmt.Printf("%-4d %-30s %-20s %-15s\n", 
			note.ID, 
			title, 
			note.CreatedAt.Format("2006-01-02 15:04"), 
			tags)
	}
}

func (nm *NoteManager) viewNote() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter note ID: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}
	
	for _, note := range nm.notes {
		if note.ID == id {
			fmt.Printf("\n=== Note %d ===\n", note.ID)
			fmt.Printf("Title: %s\n", note.Title)
			fmt.Printf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated: %s\n", note.UpdatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Tags: %s\n", strings.Join(note.Tags, ", "))
			if note.Screenshot != "" {
				fmt.Printf("Screenshot: %s\n", note.Screenshot)
			}
			fmt.Printf("\nContent:\n%s\n", note.Content)
			return
		}
	}
	
	fmt.Println("Note not found.")
}

func (nm *NoteManager) searchNotes() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter search term: ")
	term, _ := reader.ReadString('\n')
	term = strings.TrimSpace(strings.ToLower(term))
	
	if term == "" {
		fmt.Println("Search term cannot be empty.")
		return
	}
	
	var matches []Note
	for _, note := range nm.notes {
		if strings.Contains(strings.ToLower(note.Title), term) ||
		   strings.Contains(strings.ToLower(note.Content), term) ||
		   nm.containsTag(note.Tags, term) {
			matches = append(matches, note)
		}
	}
	
	if len(matches) == 0 {
		fmt.Println("No matching notes found.")
		return
	}
	
	fmt.Printf("Found %d matching notes:\n", len(matches))
	fmt.Printf("%-4s %-30s %-20s %-15s\n", "ID", "Title", "Created", "Tags")
	fmt.Println(strings.Repeat("-", 70))
	
	for _, note := range matches {
		tags := strings.Join(note.Tags, ", ")
		if len(tags) > 15 {
			tags = tags[:12] + "..."
		}
		
		title := note.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		
		fmt.Printf("%-4d %-30s %-20s %-15s\n", 
			note.ID, 
			title, 
			note.CreatedAt.Format("2006-01-02 15:04"), 
			tags)
	}
}

func (nm *NoteManager) containsTag(tags []string, term string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), term) {
			return true
		}
	}
	return false
}

func (nm *NoteManager) deleteNote() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter note ID to delete: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}
	
	for i, note := range nm.notes {
		if note.ID == id {
			fmt.Printf("Delete note '%s'? (y/n): ", note.Title)
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))
			
			if confirm == "y" {
				// Remove from memory
				nm.notes = append(nm.notes[:i], nm.notes[i+1:]...)
				
				// Remove file
				filename := filepath.Join(nm.dataDir, fmt.Sprintf("note_%d.json", note.ID))
				os.Remove(filename)
				
				// Remove screenshot if exists
				if note.Screenshot != "" {
					screenshotPath := filepath.Join(nm.dataDir, note.Screenshot)
					os.Remove(screenshotPath)
				}
				
				fmt.Println("Note deleted successfully.")
			}
			return
		}
	}
	
	fmt.Println("Note not found.")
}

func printMenu() {
	fmt.Println("\n=== THE ANCIENT SCROLLS ===")
	fmt.Println("1. Create new note")
	fmt.Println("2. List all notes")
	fmt.Println("3. View note")
	fmt.Println("4. Search notes")
	fmt.Println("5. Delete note")
	fmt.Println("6. Exit")
	fmt.Print("Choose an option: ")
}

func main() {
	fmt.Println("Welcome to The Ancient Scrolls!")
	fmt.Printf("Running on %s/%s\n", runtime.GOOS, runtime.GOARCH)
	
	nm := NewNoteManager()
	reader := bufio.NewReader(os.Stdin)
	
	for {
		printMenu()
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "1":
			nm.createNote()
		case "2":
			nm.listNotes()
		case "3":
			nm.viewNote()
		case "4":
			nm.searchNotes()
		case "5":
			nm.deleteNote()
		case "6":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
