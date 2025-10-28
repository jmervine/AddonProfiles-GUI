package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SimpleTheme is a dark theme matching WoW addon interface style
// It wraps Fyne's DarkTheme and only overrides specific colors
type SimpleTheme struct {
	base fyne.Theme
}

func NewSimpleTheme() fyne.Theme {
	return &SimpleTheme{
		base: theme.DarkTheme(),
	}
}

func (st *SimpleTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Only override specific colors, use DarkTheme for everything else
	switch name {
	case theme.ColorNamePrimary:
		// Primary button - slightly lighter gray than regular buttons
		return color.RGBA{R: 70, G: 70, B: 70, A: 255}
	case theme.ColorNameButton:
		// Regular/Cancel buttons - darker gray
		return color.RGBA{R: 50, G: 50, B: 50, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 255, G: 209, B: 0, A: 100} // Gold selection (like WoW)
	case theme.ColorNameFocus:
		return color.RGBA{R: 255, G: 209, B: 0, A: 180} // Gold focus indicator
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255} // White for links (bold)
	default:
		// Use Fyne's built-in DarkTheme for all other colors
		// This ensures proper contrast for buttons, text, etc.
		return st.base.Color(name, variant)
	}
}

func (st *SimpleTheme) Font(style fyne.TextStyle) fyne.Resource {
	return st.base.Font(style)
}

func (st *SimpleTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return st.base.Icon(name)
}

func (st *SimpleTheme) Size(name fyne.ThemeSizeName) float32 {
	return st.base.Size(name)
}
