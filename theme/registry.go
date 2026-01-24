package theme

import (
	"sort"
	"sync"
)

var (
	// mu protects the registry map
	mu sync.RWMutex
	// registry holds all registered themes
	registry = make(map[string]Palette)
	// defaultTheme is the name of the default theme
	defaultTheme = "catppuccin"
)

// Default returns the default theme palette (catppuccin).
func Default() *Palette {
	return Get(defaultTheme)
}

// Get returns a pointer to the named theme palette.
// Returns nil if the theme is not registered.
func Get(name string) *Palette {
	mu.RLock()
	defer mu.RUnlock()
	if p, ok := registry[name]; ok {
		return &p
	}
	return nil
}

// List returns a sorted list of all registered theme names.
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Register adds a theme to the registry.
// Panics if a theme with the same name is already registered.
// This function is intended to be called from init() functions in theme files.
func Register(name string, p Palette) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[name]; exists {
		panic("theme: duplicate registration of theme " + name)
	}
	registry[name] = p
}

// DefaultName returns the name of the current default theme.
// This is useful for displaying which theme is active or for config persistence.
func DefaultName() string {
	mu.RLock()
	defer mu.RUnlock()
	return defaultTheme
}

// SetDefault changes the default theme name.
// The theme must already be registered.
// Panics if the theme is not registered.
func SetDefault(name string) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[name]; !exists {
		panic("theme: cannot set default to unregistered theme " + name)
	}
	defaultTheme = name
}
