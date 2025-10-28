package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SimpleTheme is a dark theme matching WoW addon interface style
type SimpleTheme struct{}

func (st *SimpleTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 25, G: 25, B: 25, A: 255} // Dark background
	case theme.ColorNameButton:
		return color.RGBA{R: 45, G: 45, B: 45, A: 255} // Dark gray button
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 35, G: 35, B: 35, A: 255} // Darker gray
	case theme.ColorNameDisabled:
		return color.RGBA{R: 100, G: 100, B: 100, A: 255} // Dim gray for disabled text
	case theme.ColorNameForeground:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255} // Light gray/white text
	case theme.ColorNameHover:
		return color.RGBA{R: 60, G: 60, B: 60, A: 255} // Lighter gray for hover
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 35, G: 35, B: 35, A: 255} // Dark input background
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 120, G: 120, B: 120, A: 255} // Dim text for placeholders
	case theme.ColorNamePrimary:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255} // White for hyperlinks (better contrast)
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255} // White for links
	case theme.ColorNameSelection:
		return color.RGBA{R: 255, G: 209, B: 0, A: 100} // Gold selection (like WoW)
	case theme.ColorNameSuccess:
		return color.RGBA{R: 0, G: 200, B: 0, A: 255} // Green for success
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 180} // Dark shadow
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 60, G: 60, B: 60, A: 255} // Dark scrollbar
	case theme.ColorNameMenuBackground:
		return color.RGBA{R: 30, G: 30, B: 30, A: 255} // Dark menu
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 20, G: 20, B: 20, A: 240} // Dark overlay for dialogs
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (st *SimpleTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (st *SimpleTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (st *SimpleTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
