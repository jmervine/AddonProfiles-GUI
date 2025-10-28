package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	
	"github.com/jmervine/AddonProfiles-GUI/pkg/config"
	"github.com/jmervine/AddonProfiles-GUI/pkg/wow"
)

// MainWindow represents the main application window
type MainWindow struct {
	app    fyne.App
	window fyne.Window
	config *config.Config
	
	// WoW manager
	manager *wow.Manager
	
	// UI components
	profilePanel *ProfilePanel
	addonPanel   *AddonPanel
	actionPanel  *ActionPanel
	
	statusLabel *widget.Label
}

// NewMainWindow creates a new main window
func NewMainWindow(app fyne.App, cfg *config.Config) *MainWindow {
	mw := &MainWindow{
		app:         app,
		config:      cfg,
		statusLabel: widget.NewLabel("Ready"),
	}
	
	mw.window = app.NewWindow("Addon Profile Manager")
	mw.window.Resize(fyne.NewSize(1000, 600))
	mw.window.CenterOnScreen()
	
	// Check if WoW path is configured
	if cfg.WowInstallPath == "" {
		mw.showWowPathDialog()
	} else if err := cfg.Validate(); err != nil {
		mw.showWowPathDialog()
	} else {
		mw.initializeManager()
	}
	
	mw.setupUI()
	return mw
}

// showWowPathDialog shows the WoW installation path selection dialog
func (mw *MainWindow) showWowPathDialog() {
	dialog.ShowInformation("Welcome", 
		"Please select your World of Warcraft installation directory.", 
		mw.window)
	
	// Show directory picker on next event loop
	mw.window.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if e.Name == fyne.KeyEscape {
			mw.selectWowPath()
		}
	})
}

// selectWowPath shows a directory picker for WoW installation
func (mw *MainWindow) selectWowPath() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		
		if uri == nil {
			return // User cancelled
		}
		
		path := uri.Path()
		
		// Validate WoW directory
		if err := wow.ValidateWowDirectory(path); err != nil {
			dialog.ShowError(err, mw.window)
			mw.selectWowPath() // Try again
			return
		}
		
		// Save path
		mw.config.WowInstallPath = path
		if err := mw.config.Save(); err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		
		mw.initializeManager()
		mw.refresh()
		mw.setStatus("WoW directory configured")
	}, mw.window)
}

// initializeManager initializes the WoW manager
func (mw *MainWindow) initializeManager() {
	if mw.config.WowInstallPath == "" {
		return
	}
	
	// Auto-select first account if none selected
	if mw.config.SelectedAccount == "" {
		mgr := wow.NewManager(mw.config.WowInstallPath, "", 5)
		accounts, err := mgr.GetAccounts()
		if err == nil && len(accounts) > 0 {
			mw.config.SelectedAccount = accounts[0]
			mw.config.Save()
		}
	}
	
	mw.manager = wow.NewManager(
		mw.config.WowInstallPath,
		mw.config.SelectedAccount,
		mw.config.BackupCount,
	)
}

// setupUI sets up the main UI layout
func (mw *MainWindow) setupUI() {
	// Create menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Select WoW Installation", func() {
			mw.selectWowPath()
		}),
		fyne.NewMenuItem("Refresh", func() {
			mw.refresh()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Exit", func() {
			mw.app.Quit()
		}),
	)
	
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			mw.showAbout()
		}),
	)
	
	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	mw.window.SetMainMenu(mainMenu)
	
	// Create panels
	mw.profilePanel = NewProfilePanel(mw)
	mw.addonPanel = NewAddonPanel(mw)
	mw.actionPanel = NewActionPanel(mw)
	
	// Create main layout
	content := container.NewBorder(
		nil, // top
		mw.statusLabel, // bottom
		nil, // left
		nil, // right
		container.NewHSplit(
			container.NewHSplit(
				mw.profilePanel.Container(),
				mw.addonPanel.Container(),
			),
			mw.actionPanel.Container(),
		),
	)
	
	mw.window.SetContent(content)
	
	// Initial refresh
	mw.refresh()
}

// refresh refreshes all UI components
func (mw *MainWindow) refresh() {
	if mw.manager == nil {
		return
	}
	
	mw.profilePanel.Refresh()
	mw.addonPanel.Refresh()
	mw.actionPanel.Refresh()
}

// setStatus updates the status bar text
func (mw *MainWindow) setStatus(text string) {
	mw.statusLabel.SetText(text)
}

// showAbout shows the about dialog
func (mw *MainWindow) showAbout() {
	dialog.ShowInformation("About",
		"Addon Profile Manager v1.0\n\n"+
			"Manage World of Warcraft addon profiles\n"+
			"outside of the game.\n\n"+
			"github.com/jmervine/AddonProfiles-GUI",
		mw.window)
}

// ShowAndRun shows the window and runs the app
func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}

// GetManager returns the WoW manager
func (mw *MainWindow) GetManager() *wow.Manager {
	return mw.manager
}

// GetConfig returns the configuration
func (mw *MainWindow) GetConfig() *config.Config {
	return mw.config
}

// GetWindow returns the main window
func (mw *MainWindow) GetWindow() fyne.Window {
	return mw.window
}

