# Installation Security Guide

## Overview

The Ancient Scrolls installer implements several security best practices to ensure safe installation and operation on Linux systems.

## Security Features

### Installation Security

1. **No Root Requirement**
   - Installer refuses to run as root to prevent system-wide compromise
   - All files installed to user's home directory (`~/.local/`)
   - No modifications to system directories or files

2. **Path Safety**
   - Uses absolute paths and proper escaping
   - Validates all file operations
   - Creates directories with secure permissions (755)

3. **Input Validation**
   - Validates Go version compatibility
   - Checks system architecture compatibility
   - Verifies source file integrity before building

4. **Secure Build Process**
   - Builds in temporary directory
   - Cleans up temporary files after build
   - Sets appropriate executable permissions (755)

### Runtime Security

1. **File Permissions**
   - Data files created with 644 permissions (readable by user, not executable)
   - Application binary with 755 permissions (executable by user)
   - Data directory with 755 permissions

2. **Data Isolation**
   - Notes stored in user's home directory (`~/.ancient-scrolls/`)
   - No access to system files or other users' data
   - JSON format prevents code execution

3. **Input Sanitization**
   - All user inputs are properly handled
   - No shell command injection vulnerabilities
   - Screenshot commands use safe argument passing

## Installation Process Security

### Pre-Installation Checks

1. **System Compatibility**
   - Verifies Linux operating system
   - Checks processor architecture
   - Ensures Go is installed and compatible

2. **Permission Verification**
   - Ensures not running as root
   - Verifies write access to installation directories
   - Creates secure directory structure

### Installation Steps

1. **Secure Build**
   ```bash
   # Build in temporary directory
   temp_dir=$(mktemp -d)
   go build -o "$INSTALL_DIR/$APP_NAME" "$source_file"
   chmod 755 "$INSTALL_DIR/$APP_NAME"
   rm -rf "$temp_dir"
   ```

2. **Safe File Creation**
   - Desktop entries with proper metadata
   - Icons with standard permissions
   - Uninstaller with cleanup logic

3. **PATH Management**
   - Safely adds to user's PATH
   - Checks for existing entries
   - Uses shell-appropriate configuration files

## Uninstallation Security

The uninstaller:
- Only removes files it created
- Prompts before deleting user data
- Doesn't require elevated privileges
- Provides option to preserve notes

## Best Practices for Users

1. **Download Verification**
   - Only download from official repository
   - Verify source code before installation
   - Check installer script for malicious content

2. **Installation Environment**
   - Run installer from your regular user account
   - Ensure you have write access to home directory
   - Keep Go installation up to date

3. **Post-Installation**
   - Review installed files and permissions
   - Monitor data directory for unexpected changes
   - Use screenshot tools from trusted sources

## Security Considerations

### Potential Risks

1. **Screenshot Tools**
   - External screenshot utilities run with user privileges
   - Could potentially access sensitive information on screen
   - Recommendation: Only use trusted screenshot tools

2. **JSON Data**
   - Notes stored in plain text JSON format
   - Not encrypted at rest
   - Consider disk encryption for sensitive notes

3. **PATH Modification**
   - Installer modifies shell configuration files
   - Could potentially affect other applications
   - Review changes in `~/.bashrc` or `~/.zshrc`

### Mitigations

1. **Limited Scope**
   - Application only accesses its own data directory
   - No network operations or system calls
   - Minimal dependencies

2. **User Control**
   - All operations require user confirmation
   - Clear indication of what will be installed
   - Easy uninstallation process

3. **Transparent Operation**
   - Open source code for review
   - Clear documentation of all operations
   - No hidden or obfuscated functionality

## Reporting Security Issues

If you discover a security vulnerability:

1. **Do not** open a public issue
2. Contact the maintainers privately
3. Provide detailed reproduction steps
4. Allow time for investigation and patching

See SECURITY.md for complete reporting guidelines.

## Security Updates

The installer and application will be updated to address:
- Security vulnerabilities in dependencies
- New security best practices
- User-reported security concerns

Check the repository regularly for updates and security announcements.