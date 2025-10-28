package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SimpleTheme is a simple black and white theme
type SimpleTheme struct{}

func (st *SimpleTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	case theme.ColorNameButton:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255} // Light gray
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // Gray
	case theme.ColorNameForeground:
		return color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	case theme.ColorNameHover:
		return color.RGBA{R: 220, G: 220, B: 220, A: 255} // Slightly darker gray
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 122, B: 255, A: 255} // Blue
	case theme.ColorNameSelection:
		return color.RGBA{R: 255, G: 209, B: 0, A: 128} // Gold (semi-transparent)
	case theme.ColorNameSuccess:
		return color.RGBA{R: 0, G: 200, B: 0, A: 255} // Green
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

