THE ANCIENT SCROLLS: A Linux command-line note-taking app

This application takes text and screenshots, and stores them as JSON files.
It is written in the Go language and provides a secure, seamless installation experience.

## Features

- Create, view, edit, and delete notes
- Tag-based organization
- Full-text search across all notes
- Screenshot integration (with supported tools)
- JSON-based storage for portability
- Command-line interface with intuitive menu system

## Quick Installation (Recommended)

### Automatic Installer

1. **Download or clone this repository**
2. **Run the secure installer:**
   ```bash
   chmod +x install.sh
   ./install.sh
   ```

The installer will:
- ✅ Check system compatibility and Go installation
- ✅ Build the application securely
- ✅ Install to `~/.local/bin/ancient-scrolls`
- ✅ Create desktop entry for easy access
- ✅ Set up data directory with proper permissions
- ✅ Add to PATH automatically
- ✅ Create uninstaller for clean removal

### Requirements

- **Linux system** (any modern distribution)
- **Go 1.16 or higher** ([Download Go](https://golang.org/dl/))
- **Optional:** Screenshot tools (gnome-screenshot, scrot, or imagemagick)

### After Installation

Run the application:
```bash
ancient-scrolls
```

Or find "The Ancient Scrolls" in your application menu.

## Manual Installation (Advanced Users)

If you prefer to install manually:

1. **Build the application:**
   ```bash
   go build -o ancient-scrolls scrolls-init.go
   ```

2. **Move to a directory in your PATH:**
   ```bash
   mv ancient-scrolls ~/.local/bin/
   ```

3. **Make executable:**
   ```bash
   chmod +x ~/.local/bin/ancient-scrolls
   ```

## Usage

The application provides an interactive menu system:

1. **Create new note** - Add notes with title, content, tags, and optional screenshots
2. **List all notes** - View all notes in a formatted table
3. **View note** - Display full note content by ID
4. **Search notes** - Find notes by content, title, or tags
5. **Delete note** - Remove notes (with confirmation)
6. **Exit** - Close the application

### Data Storage

- Notes are stored in `~/.ancient-scrolls/` as JSON files
- Screenshots are saved in the same directory
- Each note gets a unique ID and timestamp

### Screenshot Support

For screenshot functionality, install one of:
- **Ubuntu/Debian:** `sudo apt install gnome-screenshot`
- **Fedora:** `sudo dnf install gnome-screenshot`
- **Arch:** `sudo pacman -S gnome-screenshot`
- **Alternative:** `scrot` or `imagemagick` packages

## Uninstallation

To remove the application:
```bash
ancient-scrolls-uninstall
```

This will safely remove the application while optionally preserving your notes.

## Security Features

- ✅ Refuses to run as root
- ✅ Validates all inputs and file operations
- ✅ Uses secure file permissions (755 for executables, 644 for data)
- ✅ Installs to user directory (no system-wide changes)
- ✅ Clean uninstallation process

## License

This application is freely sharable under the CC-BY-NC Creative Commons license.
https://creativecommons.org/share-your-work/cclicenses

## Support

For issues or contributions, please visit the repository at:
https://github.com/SecScholar/Note-app-project
