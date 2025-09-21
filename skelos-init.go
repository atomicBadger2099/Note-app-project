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
	Type        string    `json:"type"` // "text" or "screenshot"
	FilePath    string    `json:"file_path,omitempty"`
	Screenshot  string    `json:"screenshot,omitempty"`
}

type NotesApp struct {
	Notes      []Note `json:"notes"`
	NextID     int    `json:"next_id"`
	NotesDir   string
	ConfigFile string
}

func NewNotesApp() *NotesApp {
	homeDir, _ := os.UserHomeDir()
	notesDir := filepath.Join(homeDir, "scrolls-of-skelos")
	configFile := filepath.Join(notesDir, "scrolls.json")
	
	// Create notes directory if it doesn't exist
	os.MkdirAll(notesDir, 0755)
	os.MkdirAll(filepath.Join(notesDir, "screenshots"), 0755)
	
	app := &NotesApp{
		Notes:      []Note{},
		NextID:     1,
		NotesDir:   notesDir,
		ConfigFile: configFile,
	}
	
	app.LoadNotes()
	return app
}

func (app *NotesApp) LoadNotes() {
	if _, err := os.Stat(app.ConfigFile); os.IsNotExist(err) {
		return
	}
	
	data, err := ioutil.ReadFile(app.ConfigFile)
	if err != nil {
		fmt.Printf("Error loading notes: %v\n", err)
		return
	}
	
	if err := json.Unmarshal(data, app); err != nil {
		fmt.Printf("Error parsing notes: %v\n", err)
		return
	}
}

func (app *NotesApp) SaveNotes() {
	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling notes: %v\n", err)
		return
	}
	
	if err := ioutil.WriteFile(app.ConfigFile, data, 0644); err != nil {
		fmt.Printf("Error saving notes: %v\n", err)
	}
}

func (app *NotesApp) CreateTextNote(title, content string, tags []string) {
	note := Note{
		ID:        app.NextID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Type:      "text",
	}
	
	app.Notes = append(app.Notes, note)
	app.NextID++
	app.SaveNotes()
	
	fmt.Printf("Created scroll #%d: %s\n", note.ID, note.Title)
}

func (app *NotesApp) TakeScreenshot(title string, tags []string) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("scroll_capture_%s_%d.png", timestamp, app.NextID)
	screenshotPath := filepath.Join(app.NotesDir, "screenshots", filename)
	
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("screencapture", "-i", screenshotPath)
	case "linux":
		cmd = exec.Command("gnome-screenshot", "-a", "-f", screenshotPath)
	case "windows":
		// For Windows, we'll use a PowerShell command
		psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; Add-Type -AssemblyName System.Drawing; $Screen = [System.Windows.Forms.SystemInformation]::VirtualScreen; $Width = $Screen.Width; $Height = $Screen.Height; $Left = $Screen.Left; $Top = $Screen.Top; $bitmap = New-Object System.Drawing.Bitmap $Width, $Height; $graphic = [System.Drawing.Graphics]::FromImage($bitmap); $graphic.CopyFromScreen($Left, $Top, 0, 0, $bitmap.Size); $bitmap.Save('%s'); $graphic.Dispose(); $bitmap.Dispose()`, screenshotPath)
		cmd = exec.Command("powershell", "-Command", psScript)
	default:
		fmt.Println("Screenshot feature not supported on this platform")
		return
	}
	
	fmt.Println("Capturing ancient knowledge... (follow system prompts)")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error taking screenshot: %v\n", err)
		return
	}
	
	// Check if screenshot file was created
	if _, err := os.Stat(screenshotPath); os.IsNotExist(err) {
		fmt.Println("Knowledge capture cancelled or failed")
		return
	}
	
	note := Note{
		ID:         app.NextID,
		Title:      title,
		Tags:       tags,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Type:       "screenshot",
		FilePath:   screenshotPath,
		Screenshot: filename,
	}
	
	app.Notes = append(app.Notes, note)
	app.NextID++
	app.SaveNotes()
	
	fmt.Printf("Scroll captured and saved as scroll #%d: %s\n", note.ID, note.Title)
}

func (app *NotesApp) ListNotes() {
	if len(app.Notes) == 0 {
		fmt.Println("No scrolls found in the archives.")
		return
	}
	
	// Sort notes by creation time (newest first)
	sort.Slice(app.Notes, func(i, j int) bool {
		return app.Notes[i].CreatedAt.After(app.Notes[j].CreatedAt)
	})
	
	fmt.Println("\n=== The Ancient Scrolls ===")
	for _, note := range app.Notes {
		fmt.Printf("\n[%d] %s (%s)\n", note.ID, note.Title, note.Type)
		fmt.Printf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04"))
		if len(note.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(note.Tags, ", "))
		}
		if note.Type == "text" {
			preview := note.Content
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			fmt.Printf("Preview: %s\n", preview)
		} else {
			fmt.Printf("Captured Image: %s\n", note.Screenshot)
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}

func (app *NotesApp) ViewNote(id int) {
	for _, note := range app.Notes {
		if note.ID == id {
			fmt.Printf("\n=== Scroll of Skelos #%d ===\n", note.ID)
			fmt.Printf("Title: %s\n", note.Title)
			fmt.Printf("Type: %s\n", note.Type)
			fmt.Printf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated: %s\n", note.UpdatedAt.Format("2006-01-02 15:04:05"))
			
			if len(note.Tags) > 0 {
				fmt.Printf("Tags: %s\n", strings.Join(note.Tags, ", "))
			}
			
			if note.Type == "text" {
				fmt.Printf("\nContent:\n%s\n", note.Content)
			} else {
				fmt.Printf("\nCaptured Image: %s\n", note.Screenshot)
				fmt.Printf("File path: %s\n", note.FilePath)
				
				// Try to open the screenshot
				fmt.Print("Would you like to reveal this captured image? (y/n): ")
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				
				if response == "y" || response == "yes" {
					app.openFile(note.FilePath)
				}
			}
			return
		}
	}
	fmt.Printf("Scroll with ID %d not found in the archives.\n", id)
}

func (app *NotesApp) openFile(filePath string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	}
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error opening file: %v\n", err)
	}
}

func (app *NotesApp) SearchNotes(query string) {
	query = strings.ToLower(query)
	var matches []Note
	
	for _, note := range app.Notes {
		// Search in title, content, and tags
		if strings.Contains(strings.ToLower(note.Title), query) ||
		   strings.Contains(strings.ToLower(note.Content), query) ||
		   app.containsTag(note.Tags, query) {
			matches = append(matches, note)
		}
	}
	
	if len(matches) == 0 {
		fmt.Printf("No scrolls found containing '%s' in the archives\n", query)
		return
	}
	
	fmt.Printf("\n=== Ancient Knowledge Found: '%s' ===\n", query)
	for _, note := range matches {
		fmt.Printf("\n[%d] %s (%s)\n", note.ID, note.Title, note.Type)
		fmt.Printf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04"))
		if len(note.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(note.Tags, ", "))
		}
		if note.Type == "text" {
			preview := note.Content
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			fmt.Printf("Preview: %s\n", preview)
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}

func (app *NotesApp) containsTag(tags []string, query string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func (app *NotesApp) EditScroll(id int) {
	for i, note := range app.Notes {
		if note.ID == id {
			reader := bufio.NewReader(os.Stdin)
			
			fmt.Printf("\n=== Modifying Scroll of Skelos #%d ===\n", note.ID)
			fmt.Printf("Current Title: %s\n", note.Title)
			fmt.Printf("Type: %s\n", note.Type)
			
			if note.Type == "text" {
				// Edit text scroll
				fmt.Print("Enter new title (press Enter to keep current): ")
				newTitle, _ := reader.ReadString('\n')
				newTitle = strings.TrimSpace(newTitle)
				if newTitle != "" {
					app.Notes[i].Title = newTitle
				}
				
				fmt.Printf("Current content:\n%s\n\n", note.Content)
				fmt.Print("Enter new content (press Enter to keep current): ")
				newContent, _ := reader.ReadString('\n')
				newContent = strings.TrimSpace(newContent)
				if newContent != "" {
					app.Notes[i].Content = newContent
				}
			} else {
				// Edit image scroll title only
				fmt.Print("Enter new title (press Enter to keep current): ")
				newTitle, _ := reader.ReadString('\n')
				newTitle = strings.TrimSpace(newTitle)
				if newTitle != "" {
					app.Notes[i].Title = newTitle
				}
			}
			
			// Edit tags for both types
			if len(note.Tags) > 0 {
				fmt.Printf("Current runes (tags): %s\n", strings.Join(note.Tags, ", "))
			} else {
				fmt.Println("Current runes (tags): none")
			}
			fmt.Print("Enter new runes (comma-separated, press Enter to keep current): ")
			newTagsInput, _ := reader.ReadString('\n')
			newTagsInput = strings.TrimSpace(newTagsInput)
			
			if newTagsInput != "" {
				var newTags []string
				if newTagsInput != "" {
					newTags = strings.Split(newTagsInput, ",")
					for j, tag := range newTags {
						newTags[j] = strings.TrimSpace(tag)
					}
				}
				app.Notes[i].Tags = newTags
			}
			
			app.Notes[i].UpdatedAt = time.Now()
			app.SaveNotes()
			fmt.Printf("Scroll #%d has been modified in the archives.\n", id)
			return
		}
	}
	fmt.Printf("Scroll with ID %d not found in the archives.\n", id)
}

func (app *NotesApp) RetitleScroll(id int) {
	for i, note := range app.Notes {
		if note.ID == id {
			reader := bufio.NewReader(os.Stdin)
			
			fmt.Printf("Current title: %s\n", note.Title)
			fmt.Print("Enter new title: ")
			newTitle, _ := reader.ReadString('\n')
			newTitle = strings.TrimSpace(newTitle)
			
			if newTitle != "" {
				app.Notes[i].Title = newTitle
				app.Notes[i].UpdatedAt = time.Now()
				app.SaveNotes()
				fmt.Printf("Scroll #%d has been retitled to: %s\n", id, newTitle)
			} else {
				fmt.Println("Title unchanged.")
			}
			return
		}
	}
	fmt.Printf("Scroll with ID %d not found in the archives.\n", id)
}

func (app *NotesApp) RetagScroll(id int) {
	for i, note := range app.Notes {
		if note.ID == id {
			reader := bufio.NewReader(os.Stdin)
			
			if len(note.Tags) > 0 {
				fmt.Printf("Current runes (tags): %s\n", strings.Join(note.Tags, ", "))
			} else {
				fmt.Println("Current runes (tags): none")
			}
			
			fmt.Print("Enter new runes (comma-separated, leave empty to remove all): ")
			newTagsInput, _ := reader.ReadString('\n')
			newTagsInput = strings.TrimSpace(newTagsInput)
			
			var newTags []string
			if newTagsInput != "" {
				newTags = strings.Split(newTagsInput, ",")
				for j, tag := range newTags {
					newTags[j] = strings.TrimSpace(tag)
				}
			}
			
			app.Notes[i].Tags = newTags
			app.Notes[i].UpdatedAt = time.Now()
			app.SaveNotes()
			
			if len(newTags) > 0 {
				fmt.Printf("Scroll #%d runes updated to: %s\n", id, strings.Join(newTags, ", "))
			} else {
				fmt.Printf("All runes removed from scroll #%d\n", id)
			}
			return
		}
	}
	fmt.Printf("Scroll with ID %d not found in the archives.\n", id)
}

func (app *NotesApp) RecaptureImage(id int) {
	for i, note := range app.Notes {
		if note.ID == id {
			if note.Type != "screenshot" {
				fmt.Printf("Scroll #%d is not a captured image. Cannot recapture.\n", id)
				return
			}
			
			reader := bufio.NewReader(os.Stdin)
			
			// Ask if they want to delete the old image
			fmt.Printf("Delete the old captured image '%s'? (y/n): ", note.Screenshot)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			
			deleteOld := response == "y" || response == "yes"
			oldFilePath := note.FilePath
			
			// Create new screenshot
			timestamp := time.Now().Format("20060102_150405")
			filename := fmt.Sprintf("scroll_capture_%s_%d.png", timestamp, note.ID)
			screenshotPath := filepath.Join(app.NotesDir, "screenshots", filename)
			
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin": // macOS
				cmd = exec.Command("screencapture", "-i", screenshotPath)
			case "linux":
				cmd = exec.Command("gnome-screenshot", "-a", "-f", screenshotPath)
			case "windows":
				// For Windows, we'll use a PowerShell command
				psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; Add-Type -AssemblyName System.Drawing; $Screen = [System.Windows.Forms.SystemInformation]::VirtualScreen; $Width = $Screen.Width; $Height = $Screen.Height; $Left = $Screen.Left; $Top = $Screen.Top; $bitmap = New-Object System.Drawing.Bitmap $Width, $Height; $graphic = [System.Drawing.Graphics]::FromImage($bitmap); $graphic.CopyFromScreen($Left, $Top, 0, 0, $bitmap.Size); $bitmap.Save('%s'); $graphic.Dispose(); $bitmap.Dispose()`, screenshotPath)
				cmd = exec.Command("powershell", "-Command", psScript)
			default:
				fmt.Println("Image recapture not supported on this platform")
				return
			}
			
			fmt.Println("Recapturing ancient knowledge... (follow system prompts)")
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error recapturing image: %v\n", err)
				return
			}
			
			// Check if new screenshot file was created
			if _, err := os.Stat(screenshotPath); os.IsNotExist(err) {
				fmt.Println("Knowledge recapture cancelled or failed")
				return
			}
			
			// Update the note with new image info
			app.Notes[i].FilePath = screenshotPath
			app.Notes[i].Screenshot = filename
			app.Notes[i].UpdatedAt = time.Now()
			
			// Delete old image if requested
			if deleteOld && oldFilePath != "" {
				if err := os.Remove(oldFilePath); err != nil {
					fmt.Printf("Warning: Could not delete old image: %v\n", err)
				}
			}
			
			app.SaveNotes()
			fmt.Printf("Scroll #%d image has been recaptured: %s\n", id, filename)
			return
		}
	}
	fmt.Printf("Scroll with ID %d not found in the archives.\n", id)
}

func (app *NotesApp) DeleteNote(id int) {
	for i, note := range app.Notes {
		if note.ID == id {
			// If it's a screenshot, ask if user wants to delete the file too
			if note.Type == "screenshot" {
				fmt.Printf("Destroy the captured image '%s' from the archives as well? (y/n): ", note.Screenshot)
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				
				if response == "y" || response == "yes" {
					if err := os.Remove(note.FilePath); err != nil {
						fmt.Printf("Warning: Could not destroy captured image: %v\n", err)
					}
				}
			}
			
			// Remove note from slice
			app.Notes = append(app.Notes[:i], app.Notes[i+1:]...)
			app.SaveNotes()
			fmt.Printf("Scroll #%d has been erased from the archives.\n", id)
			return
		}
	}
	fmt.Printf("Scroll with ID %d not found in the archives.\n", id)
}

func (app *NotesApp) ShowHelp() {
	fmt.Println("\n=== The Scrolls of Skelos - Ancient Commands ===")
	fmt.Println("Available commands:")
	fmt.Println("  1 or inscribe   - Inscribe a new text scroll")
	fmt.Println("  2 or capture    - Capture an image scroll")
	fmt.Println("  3 or archive    - View all scrolls in the archive")
	fmt.Println("  4 or reveal     - Reveal a specific scroll")
	fmt.Println("  5 or seek       - Seek knowledge within scrolls")
	fmt.Println("  6 or modify     - Modify an existing scroll")
	fmt.Println("  7 or retitle    - Change a scroll's title")
	fmt.Println("  8 or retag      - Update a scroll's ancient runes")
	fmt.Println("  9 or recapture  - Replace a captured image")
	fmt.Println("  10 or erase     - Erase a scroll from existence")
	fmt.Println("  11 or wisdom    - Show these ancient commands")
	fmt.Println("  12 or depart    - Depart from the archives")
	fmt.Println()
}

func (app *NotesApp) Run() {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("🏛️  Welcome to The Scrolls of Skelos! 🏛️")
	fmt.Printf("The ancient archives are stored in: %s\n", app.NotesDir)
	app.ShowHelp()
	
	for {
		fmt.Print("\nSpeak your command, seeker of knowledge (or 'wisdom' for guidance): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		switch strings.ToLower(input) {
		case "1", "inscribe", "add":
			fmt.Print("Enter the title of your scroll: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			
			fmt.Print("Inscribe your knowledge: ")
			content, _ := reader.ReadString('\n')
			content = strings.TrimSpace(content)
			
			fmt.Print("Mark with ancient runes (tags, comma-separated, optional): ")
			tagsInput, _ := reader.ReadString('\n')
			tagsInput = strings.TrimSpace(tagsInput)
			
			var tags []string
			if tagsInput != "" {
				tags = strings.Split(tagsInput, ",")
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
			}
			
			app.CreateTextNote(title, content, tags)
			
		case "2", "capture", "screenshot":
			fmt.Print("Enter the title for your captured image: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			
			fmt.Print("Mark with ancient runes (tags, comma-separated, optional): ")
			tagsInput, _ := reader.ReadString('\n')
			tagsInput = strings.TrimSpace(tagsInput)
			
			var tags []string
			if tagsInput != "" {
				tags = strings.Split(tagsInput, ",")
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
			}
			
			app.TakeScreenshot(title, tags)
			
		case "3", "archive", "list":
			app.ListNotes()
			
		case "4", "reveal", "view":
			fmt.Print("Enter the scroll ID to reveal: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			
			if id, err := strconv.Atoi(idInput); err == nil {
				app.ViewNote(id)
			} else {
				fmt.Println("Invalid scroll ID. Please enter a number.")
			}
			
		case "5", "seek", "search":
			fmt.Print("What knowledge do you seek?: ")
			query, _ := reader.ReadString('\n')
			query = strings.TrimSpace(query)
			
			if query != "" {
				app.SearchNotes(query)
			} else {
				fmt.Println("You must speak your query to seek knowledge.")
			}
			
		case "6", "modify", "edit":
			fmt.Print("Enter the scroll ID to modify: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			
			if id, err := strconv.Atoi(idInput); err == nil {
				app.EditScroll(id)
			} else {
				fmt.Println("Invalid scroll ID. Please enter a number.")
			}
			
		case "7", "retitle":
			fmt.Print("Enter the scroll ID to retitle: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			
			if id, err := strconv.Atoi(idInput); err == nil {
				app.RetitleScroll(id)
			} else {
				fmt.Println("Invalid scroll ID. Please enter a number.")
			}
			
		case "8", "retag":
			fmt.Print("Enter the scroll ID to retag: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			
			if id, err := strconv.Atoi(idInput); err == nil {
				app.RetagScroll(id)
			} else {
				fmt.Println("Invalid scroll ID. Please enter a number.")
			}
			
		case "9", "recapture":
			fmt.Print("Enter the scroll ID to recapture: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			
			if id, err := strconv.Atoi(idInput); err == nil {
				app.RecaptureImage(id)
			} else {
				fmt.Println("Invalid scroll ID. Please enter a number.")
			}
			
		case "10", "erase", "delete":
			fmt.Print("Enter the scroll ID to erase from existence: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			
			if id, err := strconv.Atoi(idInput); err == nil {
				fmt.Printf("Are you certain you wish to erase scroll #%d from the archives? (y/n): ", id)
				confirm, _ := reader.ReadString('\n')
				confirm = strings.TrimSpace(strings.ToLower(confirm))
				
				if confirm == "y" || confirm == "yes" {
					app.DeleteNote(id)
				} else {
					fmt.Println("The scroll remains preserved in the archives.")
				}
			} else {
				fmt.Println("Invalid scroll ID. Please enter a number.")
			}
			
		case "11", "wisdom", "help":
			app.ShowHelp()
			
		case "12", "depart", "quit", "exit":
			fmt.Println("May the ancient wisdom guide you on your journey. Farewell! 🏛️")
			return
			
		default:
			fmt.Printf("Unknown command: %s\n", input)
			fmt.Println("Speak 'wisdom' to learn the ancient commands.")
		}
	}
}

func main() {
	app := NewNotesApp()
	app.Run()
}
