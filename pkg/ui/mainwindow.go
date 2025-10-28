package ui

import (
	"net/url"

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

	statusLabel  *widget.Label
	wowPathLabel *widget.Label
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
	// Show info first, then immediately open folder picker
	dialog.ShowInformation("Welcome to Addon Profile Manager",
		"This tool reads profiles from the addon's SavedVariables\n"+
			"and applies them by updating WoW's AddOns.txt file.\n\n"+
			"Click OK to select your World of Warcraft installation directory.",
		mw.window)

	// Open folder picker immediately after user closes the info dialog
	mw.selectWowPath()
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

		// Update UI
		if mw.wowPathLabel != nil {
			mw.wowPathLabel.SetText("WoW Installation: " + path)
		}

		mw.initializeManager()
		mw.refresh()
		mw.setStatus("WoW directory configured successfully")
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

	// Create header with WoW path info
	pathText := "WoW Installation: Not configured"
	if mw.config.WowInstallPath != "" {
		pathText = "WoW Installation: " + mw.config.WowInstallPath
	}
	mw.wowPathLabel = widget.NewLabel(pathText)
	mw.wowPathLabel.TextStyle = fyne.TextStyle{Bold: true}
	changePathBtn := widget.NewButton("Change WoW Path", func() {
		mw.selectWowPath()
	})

	header := container.NewBorder(
		nil,
		nil,
		mw.wowPathLabel,
		changePathBtn,
		nil,
	)

	// Create footer with status and CurseForge link
	curseforgeURL, _ := url.Parse("https://www.curseforge.com/wow/addons/addon-profiles")
	curseforgeLink := widget.NewHyperlink("▶ Get the in-game addon on CurseForge", curseforgeURL)
	curseforgeLink.TextStyle = fyne.TextStyle{Bold: true}

	footer := container.NewBorder(
		nil,
		nil,
		mw.statusLabel,
		curseforgeLink,
		nil,
	)

	// Create main layout
	content := container.NewBorder(
		header, // top
		footer, // bottom
		nil,    // left
		nil,    // right
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
