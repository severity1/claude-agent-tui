# Debug tape to capture background consistency
Output example/chatinput/screenshots/debug.gif

Set Shell "bash"
Set FontSize 14
Set Width 800
Set Height 400
Set Theme "Dracula"

Type "go run example/chatinput/main.go"
Enter
Sleep 4s

# Screenshot with placeholder
Screenshot example/chatinput/screenshots/01-placeholder.png

Sleep 2s

# Type some text
Type "hello world"
Sleep 2s

# Screenshot with text
Screenshot example/chatinput/screenshots/02-with-text.png

Sleep 1s
Ctrl+C
