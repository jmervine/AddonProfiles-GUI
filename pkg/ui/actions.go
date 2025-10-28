package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ActionPanel displays profile actions and info
type ActionPanel struct {
	mainWindow   *MainWindow
	container    *fyne.Container
	profileLabel *widget.Label
	scopeLabel   *widget.Label
	countLabel   *widget.Label
	applyBtn     *widget.Button
}

// NewActionPanel creates a new action panel
func NewActionPanel(mw *MainWindow) *ActionPanel {
	ap := &ActionPanel{
		mainWindow:   mw,
		profileLabel: widget.NewLabel("No profile selected"),
		scopeLabel:   widget.NewLabel(""),
		countLabel:   widget.NewLabel(""),
	}

	ap.applyBtn = widget.NewButton("Apply Profile", func() {
		ap.applyProfile()
	})
	ap.applyBtn.Disable()

	ap.container = container.NewVBox(
		widget.NewLabel("Settings"),
		widget.NewSeparator(),
		widget.NewLabel("Profile Name:"),
		ap.profileLabel,
		widget.NewLabel("Scope:"),
		ap.scopeLabel,
		widget.NewLabel("AddOn Count:"),
		ap.countLabel,
		widget.NewSeparator(),
		ap.applyBtn,
	)

	return ap
}

// Container returns the UI container
func (ap *ActionPanel) Container() *fyne.Container {
	return ap.container
}

// Refresh updates the action panel
func (ap *ActionPanel) Refresh() {
	profile, scope := ap.mainWindow.profilePanel.GetSelectedProfile()

	if profile == nil {
		ap.profileLabel.SetText("No profile selected")
		ap.scopeLabel.SetText("")
		ap.countLabel.SetText("")
		ap.applyBtn.Disable()
		return
	}

	ap.profileLabel.SetText(profile.Name)
	ap.scopeLabel.SetText(scope)
	ap.countLabel.SetText(fmt.Sprintf("%d addons", len(profile.Addons)))
	ap.applyBtn.Enable()
}

// applyProfile applies the selected profile
func (ap *ActionPanel) applyProfile() {
	profile, scope := ap.mainWindow.profilePanel.GetSelectedProfile()
	if profile == nil {
		return
	}

	// Confirmation dialog
	dialog.ShowConfirm(
		"Apply Profile",
		fmt.Sprintf("Apply profile '%s' (%s)?\n\nThis will update your AddOns.txt file.\nA backup will be created automatically.",
			profile.Name, scope),
		func(confirmed bool) {
			if !confirmed {
				return
			}

			mgr := ap.mainWindow.GetManager()
			if mgr == nil {
				dialog.ShowError(fmt.Errorf("WoW manager not initialized"), ap.mainWindow.GetWindow())
				return
			}

			if err := mgr.ApplyProfile(profile); err != nil {
				dialog.ShowError(err, ap.mainWindow.GetWindow())
				return
			}

			ap.mainWindow.setStatus(fmt.Sprintf("Profile '%s' applied successfully", profile.Name))
			dialog.ShowInformation("Success",
				fmt.Sprintf("Profile '%s' has been applied.\n\nYour addons will be updated when you start WoW.", profile.Name),
				ap.mainWindow.GetWindow())
		},
		ap.mainWindow.GetWindow(),
	)
}
