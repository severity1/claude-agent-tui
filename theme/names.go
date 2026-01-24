package theme

// Theme name constants for IDE autocomplete and type safety.
// Use these instead of string literals to catch typos at compile time.
//
// Example:
//
//	p := theme.Get(theme.Nord)  // IDE autocomplete works
//	p := theme.Get("nrod")      // typo not caught until runtime
const (
	Catppuccin = "catppuccin"
	Dracula    = "dracula"
	Gruvbox    = "gruvbox"
	Nord       = "nord"
	OneDark    = "onedark"
	System     = "system"
	TokyoNight = "tokyonight"
)
