# VHS tape file for ChatInput component visual regression testing
# https://github.com/charmbracelet/vhs
#
# Usage:
#   vhs validate example/chatinput/chatinput.tape  # Syntax check
#   vhs example/chatinput/chatinput.tape           # Generate screenshots
#
# This tape captures 11 key interaction states for visual regression testing.
# Includes theme switching to demonstrate the theme discovery API.

# Terminal configuration
Output example/chatinput/screenshots/chatinput-demo.gif
Set FontSize 14
Set Width 800
Set Height 600
Set Theme "Catppuccin Mocha"

# Build and run the chatinput example
Hide
Type "go run example/chatinput/main.go"
Enter
Sleep 1s
Show

# Screenshot 1: Initial state (empty input with keyhints, shows default theme description)
Sleep 500ms
Screenshot example/chatinput/screenshots/01-initial-state.png

# Screenshot 2: Theme switching - switch to dracula (key 2)
Type "2"
Sleep 500ms
Screenshot example/chatinput/screenshots/02-theme-dracula.png

# Screenshot 3: Theme switching - switch to nord (key 4)
Type "4"
Sleep 500ms
Screenshot example/chatinput/screenshots/03-theme-nord.png

# Screenshot 4: Back to catppuccin (key 1)
Type "1"
Sleep 500ms

# Screenshot 5: Typing text
Type "Hello, this is a test message"
Sleep 500ms
Screenshot example/chatinput/screenshots/04-typing-text.png

# Screenshot 6: Multi-line input (Ctrl+J for newline)
Ctrl+J
Type "This is the second line"
Sleep 500ms
Screenshot example/chatinput/screenshots/05-multi-line-input.png

# Screenshot 7: After submit (Enter)
Enter
Sleep 500ms
Screenshot example/chatinput/screenshots/06-after-submit.png

# Screenshot 8: History navigation - Up (recall submitted message)
Up
Sleep 500ms
Screenshot example/chatinput/screenshots/07-history-up.png

# Screenshot 9: History navigation - Down (return to draft/empty)
Down
Sleep 500ms
Screenshot example/chatinput/screenshots/08-history-down.png

# Screenshot 10: Focus switch to keyhints (Tab)
Tab
Sleep 500ms
Screenshot example/chatinput/screenshots/09-focus-keyhints.png

# Screenshot 11: Help window (? when hints focused)
Type "?"
Sleep 500ms
Screenshot example/chatinput/screenshots/10-help-window.png

# Screenshot 12: Clear state (Esc to close help, Esc to return to input)
Escape
Sleep 300ms
Escape
Sleep 500ms
Screenshot example/chatinput/screenshots/11-clear-state.png

# Exit the application
Ctrl+C
Sleep 500ms
