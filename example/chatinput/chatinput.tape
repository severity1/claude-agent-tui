# VHS tape file for ChatInput component visual regression testing
# https://github.com/charmbracelet/vhs
#
# Usage:
#   vhs validate example/chatinput/chatinput.tape  # Syntax check
#   vhs example/chatinput/chatinput.tape           # Generate screenshots
#
# This tape captures 17 key interaction states for visual regression testing.
# Includes theme switching, help modes, and overlay mode demonstration.

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

# Screenshot 11: Help window in Down mode (? when hints focused)
Type "?"
Sleep 500ms
Screenshot example/chatinput/screenshots/10-help-window-down.png

# Screenshot 12: Cycle to Up mode (Shift+H key)
Shift+H
Sleep 500ms
Screenshot example/chatinput/screenshots/11-help-mode-up.png

# Screenshot 13: Cycle to Top mode (Shift+H key)
Shift+H
Sleep 500ms
Screenshot example/chatinput/screenshots/12-help-mode-top.png

# Screenshot 14: Cycle to Center mode (Shift+H key)
Shift+H
Sleep 500ms
Screenshot example/chatinput/screenshots/13-help-mode-center.png

# Screenshot 15: Close help and return to Down mode for overlay test
Escape
Sleep 300ms
Shift+H
Sleep 300ms

# Screenshot 16: Enable overlay mode (Shift+O key) in Down mode
Shift+O
Sleep 500ms
Screenshot example/chatinput/screenshots/14-overlay-enabled.png

# Screenshot 17: Open help with overlay mode (? key) - shows snappier animation
Type "?"
Sleep 600ms
Screenshot example/chatinput/screenshots/15-overlay-help-open.png

# Screenshot 18: Close help with overlay mode
Escape
Sleep 600ms
Screenshot example/chatinput/screenshots/16-overlay-help-closed.png

# Screenshot 19: Disable overlay and return to normal mode
Shift+O
Sleep 300ms
Escape
Sleep 500ms
Screenshot example/chatinput/screenshots/17-clear-state.png

# Exit the application
Ctrl+C
Sleep 500ms
