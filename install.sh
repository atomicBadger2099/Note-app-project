#!/bin/bash

# The Ancient Scrolls - Secure Linux Installer
# Version: 1.0
# License: CC-BY-NC

set -euo pipefail  # Exit on error, undefined vars, pipe failures

# Constants
readonly APP_NAME="ancient-scrolls"
readonly APP_DISPLAY_NAME="The Ancient Scrolls"
readonly REPO_URL="https://github.com/SecScholar/Note-app-project"
readonly INSTALL_DIR="$HOME/.local/bin"
readonly DATA_DIR="$HOME/.ancient-scrolls"
readonly DESKTOP_DIR="$HOME/.local/share/applications"
readonly ICON_DIR="$HOME/.local/share/icons"
readonly MIN_GO_VERSION="1.16"

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root (security check)
check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "This installer should not be run as root for security reasons."
        log_info "Please run as a regular user. The installer will install to your home directory."
        exit 1
    fi
}

# Check system compatibility
check_system() {
    log_info "Checking system compatibility..."
    
    if [[ "$OSTYPE" != "linux-gnu"* ]]; then
        log_error "This installer is designed for Linux systems only."
        log_info "Detected OS: $OSTYPE"
        exit 1
    fi
    
    # Check architecture
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64|amd64)
            log_success "Architecture: $arch (supported)"
            ;;
        *)
            log_warning "Architecture: $arch (may not be supported)"
            log_info "The application should work but hasn't been tested on this architecture."
            read -p "Continue anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
            ;;
    esac
}

# Check Go installation
check_go() {
    log_info "Checking Go installation..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed on your system."
        log_info "Please install Go from https://golang.org/dl/"
        log_info "Minimum required version: $MIN_GO_VERSION"
        
        # Offer to help with installation
        echo
        read -p "Would you like to see Go installation instructions? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            show_go_install_instructions
        fi
        exit 1
    fi
    
    # Check Go version
    local go_version
    go_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+(\.[0-9]+)?')
    
    if ! version_compare "$go_version" "$MIN_GO_VERSION"; then
        log_error "Go version $go_version is too old."
        log_info "Minimum required version: $MIN_GO_VERSION"
        log_info "Please update Go from https://golang.org/dl/"
        exit 1
    fi
    
    log_success "Go $go_version is installed and compatible"
}

# Compare version numbers
version_compare() {
    local version1=$1
    local min_version=$2
    
    # Convert versions to comparable format
    local v1_major v1_minor v1_patch
    local v2_major v2_minor v2_patch
    
    IFS='.' read -r v1_major v1_minor v1_patch <<< "$version1.0.0"
    IFS='.' read -r v2_major v2_minor v2_patch <<< "$min_version.0.0"
    
    # Compare major version
    if [[ $v1_major -gt $v2_major ]]; then
        return 0
    elif [[ $v1_major -lt $v2_major ]]; then
        return 1
    fi
    
    # Compare minor version
    if [[ $v1_minor -gt $v2_minor ]]; then
        return 0
    elif [[ $v1_minor -lt $v2_minor ]]; then
        return 1
    fi
    
    # Compare patch version
    if [[ $v1_patch -ge $v2_patch ]]; then
        return 0
    else
        return 1
    fi
}

# Show Go installation instructions
show_go_install_instructions() {
    cat << 'EOF'

=== Go Installation Instructions ===

1. Download Go from https://golang.org/dl/
2. For most Linux systems, run:
   
   # Download (replace with latest version)
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   
   # Remove old installation and extract new one
   sudo rm -rf /usr/local/go
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   
   # Add to PATH (add to ~/.bashrc or ~/.profile)
   export PATH=$PATH:/usr/local/go/bin

3. Restart your terminal or run: source ~/.bashrc
4. Verify with: go version

EOF
}

# Create necessary directories
create_directories() {
    log_info "Creating directories..."
    
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$DATA_DIR"
    mkdir -p "$DESKTOP_DIR"
    mkdir -p "$ICON_DIR"
    
    # Set secure permissions
    chmod 755 "$INSTALL_DIR"
    chmod 755 "$DATA_DIR"
    chmod 755 "$DESKTOP_DIR"
    chmod 755 "$ICON_DIR"
    
    log_success "Directories created"
}

# Build the application
build_application() {
    log_info "Building The Ancient Scrolls..."
    
    local temp_dir
    temp_dir=$(mktemp -d)
    local source_file="$temp_dir/scrolls-init.go"
    
    # Copy source file to temp directory
    cp "scrolls-init.go" "$source_file"
    
    # Build the application
    if ! go build -o "$INSTALL_DIR/$APP_NAME" "$source_file"; then
        log_error "Failed to build the application"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Set executable permissions
    chmod 755 "$INSTALL_DIR/$APP_NAME"
    
    # Clean up
    rm -rf "$temp_dir"
    
    log_success "Application built successfully"
}

# Create desktop entry
create_desktop_entry() {
    log_info "Creating desktop entry..."
    
    local desktop_file="$DESKTOP_DIR/$APP_NAME.desktop"
    
    cat > "$desktop_file" << EOF
[Desktop Entry]
Name=$APP_DISPLAY_NAME
Comment=A command-line note-taking application with screenshot support
Exec=x-terminal-emulator -e $INSTALL_DIR/$APP_NAME
Icon=$APP_NAME
Terminal=true
Type=Application
Categories=Office;TextEditor;Utility;
Keywords=notes;text;editor;screenshots;
StartupNotify=false
EOF
    
    chmod 644 "$desktop_file"
    
    # Update desktop database if available
    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi
    
    log_success "Desktop entry created"
}

# Create a simple icon
create_icon() {
    log_info "Creating application icon..."
    
    local icon_file="$ICON_DIR/$APP_NAME.svg"
    
    # Create a simple SVG icon
    cat > "$icon_file" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<svg width="48" height="48" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
  <rect width="48" height="48" rx="8" fill="#2c3e50"/>
  <text x="24" y="30" font-family="serif" font-size="24" font-weight="bold" text-anchor="middle" fill="#ecf0f1">📜</text>
</svg>
EOF
    
    chmod 644 "$icon_file"
    log_success "Icon created"
}

# Add to PATH if needed
update_path() {
    local shell_rc="$HOME/.bashrc"
    
    # Detect shell and use appropriate rc file
    case "$SHELL" in
        */zsh)
            shell_rc="$HOME/.zshrc"
            ;;
        */fish)
            shell_rc="$HOME/.config/fish/config.fish"
            return  # Fish uses different syntax, skip for now
            ;;
    esac
    
    # Check if INSTALL_DIR is already in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        log_info "Adding $INSTALL_DIR to PATH..."
        
        echo "" >> "$shell_rc"
        echo "# Added by The Ancient Scrolls installer" >> "$shell_rc"
        echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$shell_rc"
        
        log_success "Added to PATH in $shell_rc"
        log_info "Please restart your terminal or run: source $shell_rc"
    else
        log_info "$INSTALL_DIR is already in PATH"
    fi
}

# Check for optional screenshot dependencies
check_screenshot_deps() {
    log_info "Checking screenshot dependencies..."
    
    local tools=("gnome-screenshot" "scrot" "import")
    local found=false
    
    for tool in "${tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            log_success "Found screenshot tool: $tool"
            found=true
        fi
    done
    
    if ! $found; then
        log_warning "No screenshot tools found"
        log_info "Install one of the following for screenshot support:"
        log_info "  - gnome-screenshot (GNOME)"
        log_info "  - scrot (lightweight)"
        log_info "  - imagemagick (for 'import' command)"
        echo
        log_info "Example installations:"
        log_info "  Ubuntu/Debian: sudo apt install gnome-screenshot"
        log_info "  Fedora: sudo dnf install gnome-screenshot"
        log_info "  Arch: sudo pacman -S gnome-screenshot"
    fi
}

# Create uninstaller
create_uninstaller() {
    log_info "Creating uninstaller..."
    
    local uninstall_script="$INSTALL_DIR/${APP_NAME}-uninstall"
    
    cat > "$uninstall_script" << EOF
#!/bin/bash

# The Ancient Scrolls Uninstaller

echo "Uninstalling The Ancient Scrolls..."

# Remove binary
rm -f "$INSTALL_DIR/$APP_NAME"

# Remove desktop entry
rm -f "$DESKTOP_DIR/$APP_NAME.desktop"

# Remove icon
rm -f "$ICON_DIR/$APP_NAME.svg"

# Remove this uninstaller
rm -f "$uninstall_script"

# Ask about data directory
echo
read -p "Remove data directory $DATA_DIR? This will delete all your notes! (y/N): " -n 1 -r
echo
if [[ \$REPLY =~ ^[Yy]$ ]]; then
    rm -rf "$DATA_DIR"
    echo "Data directory removed."
else
    echo "Data directory preserved."
fi

echo "The Ancient Scrolls has been uninstalled."
echo "You may need to remove the PATH entry from your shell configuration manually."
EOF
    
    chmod 755 "$uninstall_script"
    log_success "Uninstaller created at $uninstall_script"
}

# Test installation
test_installation() {
    log_info "Testing installation..."
    
    if [[ -x "$INSTALL_DIR/$APP_NAME" ]]; then
        log_success "Installation test passed"
        return 0
    else
        log_error "Installation test failed"
        return 1
    fi
}

# Main installation function
main() {
    echo "======================================"
    echo "  The Ancient Scrolls - Linux Installer"
    echo "======================================"
    echo
    
    log_info "Starting secure installation..."
    
    # Security and compatibility checks
    check_root
    check_system
    check_go
    
    # Verify we have the source file
    if [[ ! -f "scrolls-init.go" ]]; then
        log_error "scrolls-init.go not found in current directory"
        log_info "Please run this installer from the directory containing scrolls-init.go"
        exit 1
    fi
    
    echo
    log_info "Installation will proceed with the following:"
    log_info "  Application: $INSTALL_DIR/$APP_NAME"
    log_info "  Data directory: $DATA_DIR"
    log_info "  Desktop entry: $DESKTOP_DIR/$APP_NAME.desktop"
    echo
    
    read -p "Continue with installation? (Y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        log_info "Installation cancelled by user"
        exit 0
    fi
    
    # Perform installation
    create_directories
    build_application
    create_desktop_entry
    create_icon
    update_path
    create_uninstaller
    
    # Additional checks
    check_screenshot_deps
    
    # Test installation
    if test_installation; then
        echo
        echo "======================================"
        log_success "Installation completed successfully!"
        echo "======================================"
        echo
        log_info "You can now run the application with: $APP_NAME"
        log_info "Or use the desktop entry: $APP_DISPLAY_NAME"
        log_info "Data will be stored in: $DATA_DIR"
        echo
        log_info "To uninstall, run: ${APP_NAME}-uninstall"
        echo
        
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            log_warning "Please restart your terminal or run: source ~/.bashrc"
        fi
    else
        log_error "Installation failed"
        exit 1
    fi
}

# Run installer
main "$@"