// Package theme provides a centralized theming system for claude-agent-tui components.
//
// This package defines a Palette struct with 50 semantic color definitions organized
// into four categories: Core UI (15), Diff (12), Markdown (14), and Syntax (9).
// All colors use lipgloss.AdaptiveColor to automatically adjust for light/dark terminals.
//
// # Built-in Themes
//
// Seven themes are available out of the box:
//
//   - catppuccin: Soothing pastel theme (default) - Mocha/Latte variants
//   - dracula: Classic dark theme with vibrant colors
//   - gruvbox: Retro groove with warm earthy tones
//   - nord: Arctic north-bluish color palette
//   - tokyonight: Clean dark theme inspired by Tokyo city lights
//   - onedark: The iconic Atom editor theme
//   - system: ANSI 16-color fallback for maximum compatibility
//
// # Usage
//
// Use the default theme:
//
//	p := theme.Default()
//	input := chatinput.New(chatinput.WithPalette(p))
//
// Use a specific theme with IDE autocomplete:
//
//	p := theme.Get(theme.Dracula)  // constants provide autocomplete
//	input := chatinput.New(chatinput.WithPalette(p))
//
// Get the default theme name:
//
//	name := theme.DefaultName()  // returns "catppuccin"
//
// List available themes:
//
//	for _, name := range theme.List() {
//	    fmt.Println(name)
//	}
//
// List themes with descriptions for UI selection menus:
//
//	for _, info := range theme.ListInfo() {
//	    fmt.Printf("%s - %s\n", info.Name, info.Description)
//	}
//	// Output:
//	//   catppuccin - Soothing pastel theme (Mocha/Latte)
//	//   dracula - Classic dark theme with vibrant colors
//	//   ...
//
// Get info for a specific theme:
//
//	if info := theme.GetInfo(theme.Nord); info != nil {
//	    fmt.Printf("Using %s: %s\n", info.Name, info.Description)
//	}
//
// # Custom Themes
//
// Register a custom theme:
//
//	theme.Register("custom", theme.Palette{
//	    Primary: lipgloss.AdaptiveColor{Light: "#000", Dark: "#fff"},
//	    // ... all 45 fields
//	})
//
// # Component Integration
//
// Components accept palettes via the WithPalette option:
//
//	// ChatInput
//	input := chatinput.New(chatinput.WithPalette(theme.Get(theme.Nord)))
//
//	// KeyHints
//	hints := keyhints.New(bindings, keyhints.WithPalette(theme.Get(theme.Nord)))
//
// Components also provide StylesFromPalette() for fine-grained control:
//
//	p := theme.Get(theme.Gruvbox)
//	styles := chatinput.StylesFromPalette(p)
//	styles.Border = lipgloss.RoundedBorder()  // Override just the border
//	input := chatinput.New(chatinput.WithStyles(styles))
//
// # Importing Themes
//
// To use built-in themes, import the themes package for its side effects:
//
//	import _ "github.com/severity1/claude-agent-tui/theme/themes"
//
// This registers all built-in themes via init() functions. Without this import,
// only manually registered themes will be available.
package theme
