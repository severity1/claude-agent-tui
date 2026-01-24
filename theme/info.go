package theme

import "sort"

// Info contains metadata about a theme for display in UI selection menus.
type Info struct {
	Name        string // Theme identifier (e.g., "catppuccin")
	Description string // Human-readable description for UI display
}

// themeDescriptions maps theme names to their descriptions.
var themeDescriptions = map[string]string{
	Catppuccin: "Soothing pastel theme (Mocha/Latte)",
	Dracula:    "Classic dark theme with vibrant colors",
	Gruvbox:    "Retro groove with warm earthy tones",
	Nord:       "Arctic north-bluish color palette",
	OneDark:    "Atom editor's iconic dark theme",
	System:     "ANSI 16-color for max compatibility",
	TokyoNight: "Clean theme inspired by Tokyo city lights",
}

// GetInfo returns metadata for the named theme.
// Returns nil if the theme is not registered.
func GetInfo(name string) *Info {
	mu.RLock()
	defer mu.RUnlock()

	if _, ok := registry[name]; !ok {
		return nil
	}

	desc, ok := themeDescriptions[name]
	if !ok {
		desc = ""
	}

	return &Info{
		Name:        name,
		Description: desc,
	}
}

// ListInfo returns metadata for all registered themes, sorted by name.
// This is useful for building theme selection UIs.
//
// Example:
//
//	for _, info := range theme.ListInfo() {
//	    fmt.Printf("%s - %s\n", info.Name, info.Description)
//	}
func ListInfo() []Info {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)

	infos := make([]Info, len(names))
	for i, name := range names {
		desc, ok := themeDescriptions[name]
		if !ok {
			desc = ""
		}
		infos[i] = Info{
			Name:        name,
			Description: desc,
		}
	}

	return infos
}
