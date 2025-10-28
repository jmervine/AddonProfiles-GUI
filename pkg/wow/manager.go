package wow

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmervine/AddonProfiles-GUI/pkg/lua"
)

// Manager handles WoW data operations
type Manager struct {
	wowPath         string
	selectedAccount string
	backupCount     int
}

// NewManager creates a new WoW data manager
func NewManager(wowPath, account string, backupCount int) *Manager {
	return &Manager{
		wowPath:         wowPath,
		selectedAccount: account,
		backupCount:     backupCount,
	}
}

// GetAccounts returns a list of account names found in the WTF directory
func (m *Manager) GetAccounts() ([]string, error) {
	accountsPath := filepath.Join(m.wowPath, "WTF", "Account")

	entries, err := os.ReadDir(accountsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts directory: %w", err)
	}

	var accounts []string
	for _, entry := range entries {
		if entry.IsDir() {
			accounts = append(accounts, entry.Name())
		}
	}

	return accounts, nil
}

// LoadProfiles loads all profiles from SavedVariables
func (m *Manager) LoadProfiles() (*lua.Database, error) {
	if m.selectedAccount == "" {
		return nil, fmt.Errorf("no account selected")
	}

	savedVarsPath := filepath.Join(m.wowPath, "WTF", "Account", m.selectedAccount,
		"SavedVariables", "AddonProfilesDB.lua")

	if _, err := os.Stat(savedVarsPath); os.IsNotExist(err) {
		// Return empty database if file doesn't exist
		return &lua.Database{
			Global: struct {
				ActiveProfile string
				Profiles      map[string]*lua.Profile
				Settings      map[string]interface{}
			}{
				Profiles: make(map[string]*lua.Profile),
				Settings: make(map[string]interface{}),
			},
			Char: make(map[string]struct {
				ActiveProfile string
				Profiles      map[string]*lua.Profile
			}),
		}, nil
	}

	return lua.ParseFile(savedVarsPath)
}

// GetActiveAddons returns the currently active addons from AddOns.txt
func (m *Manager) GetActiveAddons() (map[string]bool, error) {
	if m.selectedAccount == "" {
		return nil, fmt.Errorf("no account selected")
	}

	addonsPath := filepath.Join(m.wowPath, "WTF", "Account", m.selectedAccount, "AddOns.txt")

	if _, err := os.Stat(addonsPath); os.IsNotExist(err) {
		return make(map[string]bool), nil
	}

	return parseAddOnsFile(addonsPath)
}

// ApplyProfile applies a profile by updating AddOns.txt
func (m *Manager) ApplyProfile(profile *lua.Profile) error {
	if m.selectedAccount == "" {
		return fmt.Errorf("no account selected")
	}

	addonsPath := filepath.Join(m.wowPath, "WTF", "Account", m.selectedAccount, "AddOns.txt")

	// Create backup
	if err := m.createBackup(addonsPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Write new AddOns.txt
	if err := writeAddOnsFile(addonsPath, profile.Addons); err != nil {
		return fmt.Errorf("failed to write AddOns.txt: %w", err)
	}

	// Clean up old backups
	if err := m.cleanupBackups(addonsPath); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to cleanup old backups: %v\n", err)
	}

	return nil
}

// createBackup creates a timestamped backup of AddOns.txt
func (m *Manager) createBackup(addonsPath string) error {
	if _, err := os.Stat(addonsPath); os.IsNotExist(err) {
		// No file to backup
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup.%s", addonsPath, timestamp)

	data, err := os.ReadFile(addonsPath)
	if err != nil {
		return fmt.Errorf("failed to read AddOns.txt: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}

// cleanupBackups removes old backups, keeping only the most recent N
func (m *Manager) cleanupBackups(addonsPath string) error {
	dir := filepath.Dir(addonsPath)
	base := filepath.Base(addonsPath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Find all backup files
	var backups []string
	prefix := base + ".backup."
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			backups = append(backups, filepath.Join(dir, entry.Name()))
		}
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		infoI, _ := os.Stat(backups[i])
		infoJ, _ := os.Stat(backups[j])
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Remove old backups
	if len(backups) > m.backupCount {
		for _, backup := range backups[m.backupCount:] {
			os.Remove(backup)
		}
	}

	return nil
}

// parseAddOnsFile parses an AddOns.txt file
func parseAddOnsFile(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	addons := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Format: "AddonName: 1" (enabled) or "# AddonName: 0" (disabled)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		addonName := strings.TrimSpace(parts[0])
		enabled := strings.TrimSpace(parts[1]) == "1"

		// Remove comment marker if present
		if strings.HasPrefix(addonName, "#") {
			addonName = strings.TrimSpace(strings.TrimPrefix(addonName, "#"))
			enabled = false
		}

		addons[addonName] = enabled
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return addons, nil
}

// writeAddOnsFile writes addons to AddOns.txt file
func writeAddOnsFile(path string, addons map[string]bool) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Sort addon names for consistent output
	var names []string
	for name := range addons {
		names = append(names, name)
	}
	sort.Strings(names)

	// Write each addon
	for _, name := range names {
		enabled := addons[name]
		if enabled {
			fmt.Fprintf(writer, "%s: 1\n", name)
		} else {
			fmt.Fprintf(writer, "# %s: 0\n", name)
		}
	}

	return writer.Flush()
}

// ValidateWowDirectory checks if a directory is a valid WoW installation
func ValidateWowDirectory(path string) error {
	// Check if directory exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("directory does not exist: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Check for WTF directory
	wtfPath := filepath.Join(path, "WTF")
	if _, err := os.Stat(wtfPath); os.IsNotExist(err) {
		return fmt.Errorf("WTF directory not found")
	}

	// Check for Account directory
	accountPath := filepath.Join(wtfPath, "Account")
	if _, err := os.Stat(accountPath); os.IsNotExist(err) {
		return fmt.Errorf("WTF/Account directory not found")
	}

	return nil
}
