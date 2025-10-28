package ui

import (
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AddonPanel displays the addon list for selected profile
type AddonPanel struct {
	mainWindow  *MainWindow
	container   *fyne.Container
	addonList   *widget.List
	searchEntry *widget.Entry
	addons      []AddonItem
	filtered    []AddonItem
}

// AddonItem represents an addon
type AddonItem struct {
	Name    string
	Enabled bool
}

// NewAddonPanel creates a new addon panel
func NewAddonPanel(mw *MainWindow) *AddonPanel {
	ap := &AddonPanel{
		mainWindow: mw,
		addons:     []AddonItem{},
		filtered:   []AddonItem{},
	}
	
	ap.searchEntry = widget.NewEntry()
	ap.searchEntry.SetPlaceHolder("Search addons...")
	ap.searchEntry.OnChanged = func(text string) {
		ap.filterAddons(text)
	}
	
	ap.addonList = widget.NewList(
		func() int {
			return len(ap.filtered)
		},
		func() fyne.CanvasObject {
			return widget.NewCheck("Addon Name", nil)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			check := obj.(*widget.Check)
			if id < len(ap.filtered) {
				addon := ap.filtered[id]
				check.SetText(addon.Name)
				check.SetChecked(addon.Enabled)
				check.Disable() // Read-only
			}
		},
	)
	
	ap.container = container.NewBorder(
		container.NewVBox(
			widget.NewLabel("AddOns"),
			ap.searchEntry,
		),
		nil,
		nil,
		nil,
		ap.addonList,
	)
	
	return ap
}

// Container returns the UI container
func (ap *AddonPanel) Container() *fyne.Container {
	return ap.container
}

// Refresh reloads the addon list
func (ap *AddonPanel) Refresh() {
	profile, _ := ap.mainWindow.profilePanel.GetSelectedProfile()
	if profile == nil {
		ap.addons = []AddonItem{}
		ap.filtered = []AddonItem{}
		ap.addonList.Refresh()
		return
	}
	
	// Convert map to sorted slice
	ap.addons = []AddonItem{}
	for name, enabled := range profile.Addons {
		ap.addons = append(ap.addons, AddonItem{
			Name:    name,
			Enabled: enabled,
		})
	}
	
	sort.Slice(ap.addons, func(i, j int) bool {
		return ap.addons[i].Name < ap.addons[j].Name
	})
	
	ap.filterAddons(ap.searchEntry.Text)
}

// filterAddons filters the addon list based on search text
func (ap *AddonPanel) filterAddons(search string) {
	if search == "" {
		ap.filtered = ap.addons
	} else {
		ap.filtered = []AddonItem{}
		searchLower := strings.ToLower(search)
		for _, addon := range ap.addons {
			if strings.Contains(strings.ToLower(addon.Name), searchLower) {
				ap.filtered = append(ap.filtered, addon)
			}
		}
	}
	
	ap.addonList.Refresh()
}

