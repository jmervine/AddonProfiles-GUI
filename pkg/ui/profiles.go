package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	
	"github.com/jmervine/AddonProfiles-GUI/pkg/lua"
)

// ProfilePanel displays the profile list
type ProfilePanel struct {
	mainWindow      *MainWindow
	container       *fyne.Container
	profileList     *widget.List
	selectedProfile *lua.Profile
	selectedScope   string
	profiles        []*ProfileItem
}

// ProfileItem represents a profile in the list
type ProfileItem struct {
	Name     string
	Scope    string
	IsActive bool
	Profile  *lua.Profile
}

// NewProfilePanel creates a new profile panel
func NewProfilePanel(mw *MainWindow) *ProfilePanel {
	pp := &ProfilePanel{
		mainWindow: mw,
		profiles:   []*ProfileItem{},
	}
	
	pp.profileList = widget.NewList(
		func() int {
			return len(pp.profiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Profile Name")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(pp.profiles) {
				item := pp.profiles[id]
				prefix := ""
				if item.IsActive {
					prefix = "* "
				}
				label.SetText(fmt.Sprintf("%s%s (%s)", prefix, item.Name, item.Scope))
			}
		},
	)
	
	pp.profileList.OnSelected = func(id widget.ListItemID) {
		if id < len(pp.profiles) {
			pp.selectedProfile = pp.profiles[id].Profile
			pp.selectedScope = pp.profiles[id].Scope
			mw.addonPanel.Refresh()
			mw.actionPanel.Refresh()
		}
	}
	
	refreshBtn := widget.NewButton("Refresh Profiles", func() {
		pp.Refresh()
	})
	
	pp.container = container.NewBorder(
		widget.NewLabel("Profiles"),
		refreshBtn,
		nil,
		nil,
		pp.profileList,
	)
	
	return pp
}

// Container returns the UI container
func (pp *ProfilePanel) Container() *fyne.Container {
	return pp.container
}

// Refresh reloads the profile list
func (pp *ProfilePanel) Refresh() {
	mgr := pp.mainWindow.GetManager()
	if mgr == nil {
		return
	}
	
	db, err := mgr.LoadProfiles()
	if err != nil {
		pp.mainWindow.setStatus(fmt.Sprintf("Error loading profiles: %v", err))
		return
	}
	
	pp.profiles = []*ProfileItem{}
	
	// Add global profiles
	for name, profile := range db.Global.Profiles {
		pp.profiles = append(pp.profiles, &ProfileItem{
			Name:     name,
			Scope:    "account",
			IsActive: name == db.Global.ActiveProfile,
			Profile:  profile,
		})
	}
	
	// Add character profiles (current character only for simplicity)
	// You could expand this to show all characters
	
	pp.profileList.Refresh()
	pp.mainWindow.setStatus(fmt.Sprintf("Loaded %d profiles", len(pp.profiles)))
}

// GetSelectedProfile returns the currently selected profile
func (pp *ProfilePanel) GetSelectedProfile() (*lua.Profile, string) {
	return pp.selectedProfile, pp.selectedScope
}

