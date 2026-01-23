# VHS tape file for ChatInput component visual regression testing
# https://github.com/charmbracelet/vhs
#
# Usage:
#   vhs validate example/chatinput/chatinput.tape  # Syntax check
#   vhs example/chatinput/chatinput.tape           # Generate screenshots
#
# This tape captures 9 key interaction states for visual regression testing.

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

# Screenshot 1: Initial state (empty input with keyhints)
Sleep 500ms
Screenshot example/chatinput/screenshots/01-initial-state.png

# Screenshot 2: Typing text
Type "Hello, this is a test message"
Sleep 500ms
Screenshot example/chatinput/screenshots/02-typing-text.png

# Screenshot 3: Multi-line input (Ctrl+J for newline)
Ctrl+J
Type "This is the second line"
Sleep 500ms
Screenshot example/chatinput/screenshots/03-multi-line-input.png

# Screenshot 4: After submit (Enter)
Enter
Sleep 500ms
Screenshot example/chatinput/screenshots/04-after-submit.png

# Screenshot 5: History navigation - Up (recall submitted message)
Up
Sleep 500ms
Screenshot example/chatinput/screenshots/05-history-up.png

# Screenshot 6: History navigation - Down (return to draft/empty)
Down
Sleep 500ms
Screenshot example/chatinput/screenshots/06-history-down.png

# Screenshot 7: Focus switch to keyhints (Tab)
Tab
Sleep 500ms
Screenshot example/chatinput/screenshots/07-focus-keyhints.png

# Screenshot 8: Help window (? when hints focused)
Type "?"
Sleep 500ms
Screenshot example/chatinput/screenshots/08-help-window.png

# Screenshot 9: Clear state (Esc to close help, Esc to return to input)
Escape
Sleep 300ms
Escape
Sleep 500ms
Screenshot example/chatinput/screenshots/09-clear-state.png

# Exit the application
Ctrl+C
Sleep 500ms
