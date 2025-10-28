package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/jmervine/AddonProfiles-GUI/pkg/config"
	"github.com/jmervine/AddonProfiles-GUI/pkg/ui"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Fyne application
	myApp := app.NewWithID("com.github.jmervine.addonprofiles")
	myApp.Settings().SetTheme(&ui.SimpleTheme{})

	// Create main window
	mainWindow := ui.NewMainWindow(myApp, cfg)
	
	// Show and run
	mainWindow.ShowAndRun()
}

