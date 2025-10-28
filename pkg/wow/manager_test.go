package wow

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmervine/AddonProfiles-GUI/pkg/lua"
)

func TestNewManager(t *testing.T) {
	mgr := NewManager("/path/to/wow", "TestAccount", 5)

	if mgr == nil {
		t.Fatal("NewManager returned nil")
	}

	if mgr.wowPath != "/path/to/wow" {
		t.Errorf("wowPath = %v, want /path/to/wow", mgr.wowPath)
	}

	if mgr.selectedAccount != "TestAccount" {
		t.Errorf("selectedAccount = %v, want TestAccount", mgr.selectedAccount)
	}

	if mgr.backupCount != 5 {
		t.Errorf("backupCount = %v, want 5", mgr.backupCount)
	}
}

func TestParseAddOnsFile(t *testing.T) {
	path := filepath.Join("testdata", "AddOns.txt")
	addons, err := parseAddOnsFile(path)
	if err != nil {
		t.Fatalf("parseAddOnsFile() error = %v", err)
	}

	expected := map[string]bool{
		"Ace3":          true,
		"AddonProfiles": false,
		"Details":       true,
		"DBM-Core":      true,
		"WeakAuras":     false,
		"BigWigs":       true,
	}

	if len(addons) != len(expected) {
		t.Errorf("Expected %d addons, got %d", len(expected), len(addons))
	}

	for name, enabled := range expected {
		if addons[name] != enabled {
			t.Errorf("Addon %s: got %v, want %v", name, addons[name], enabled)
		}
	}
}

func TestWriteAddOnsFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "AddOns.txt")

	addons := map[string]bool{
		"Addon1": true,
		"Addon2": false,
		"Addon3": true,
	}

	if err := writeAddOnsFile(testPath, addons); err != nil {
		t.Fatalf("writeAddOnsFile() error = %v", err)
	}

	// Read back and verify
	parsed, err := parseAddOnsFile(testPath)
	if err != nil {
		t.Fatalf("parseAddOnsFile() error = %v", err)
	}

	if len(parsed) != len(addons) {
		t.Errorf("Expected %d addons, got %d", len(addons), len(parsed))
	}

	for name, enabled := range addons {
		if parsed[name] != enabled {
			t.Errorf("Addon %s: got %v, want %v", name, parsed[name], enabled)
		}
	}
}

func TestValidateWowDirectory(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
		cleanup func(string)
	}{
		{
			name: "non-existent directory",
			setup: func() string {
				return "/non/existent/path"
			},
			wantErr: true,
		},
		{
			name: "valid WoW directory",
			setup: func() string {
				tmpDir, _ := os.MkdirTemp("", "wow-test-*")
				os.MkdirAll(filepath.Join(tmpDir, "WTF", "Account"), 0755)
				return tmpDir
			},
			wantErr: false,
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
		{
			name: "directory without WTF",
			setup: func() string {
				tmpDir, _ := os.MkdirTemp("", "wow-test-*")
				return tmpDir
			},
			wantErr: true,
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
		{
			name: "directory with WTF but no Account",
			setup: func() string {
				tmpDir, _ := os.MkdirTemp("", "wow-test-*")
				os.Mkdir(filepath.Join(tmpDir, "WTF"), 0755)
				return tmpDir
			},
			wantErr: true,
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup()
			if tt.cleanup != nil {
				defer tt.cleanup(dir)
			}

			err := ValidateWowDirectory(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWowDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAccounts(t *testing.T) {
	// Create test structure
	tmpDir, err := os.MkdirTemp("", "wow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	accountsDir := filepath.Join(tmpDir, "WTF", "Account")
	os.MkdirAll(accountsDir, 0755)

	// Create test accounts
	accounts := []string{"Account1", "Account2", "Account3"}
	for _, acc := range accounts {
		os.Mkdir(filepath.Join(accountsDir, acc), 0755)
	}

	// Create a file (should be ignored)
	os.WriteFile(filepath.Join(accountsDir, "file.txt"), []byte("test"), 0644)

	mgr := NewManager(tmpDir, "", 5)
	found, err := mgr.GetAccounts()
	if err != nil {
		t.Fatalf("GetAccounts() error = %v", err)
	}

	if len(found) != len(accounts) {
		t.Errorf("Expected %d accounts, got %d", len(accounts), len(found))
	}

	for _, acc := range accounts {
		foundIt := false
		for _, f := range found {
			if f == acc {
				foundIt = true
				break
			}
		}
		if !foundIt {
			t.Errorf("Account %s not found", acc)
		}
	}
}

func TestLoadProfiles(t *testing.T) {
	// Create test structure
	tmpDir, err := os.MkdirTemp("", "wow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	account := "TestAccount"
	savedVarsDir := filepath.Join(tmpDir, "WTF", "Account", account, "SavedVariables")
	os.MkdirAll(savedVarsDir, 0755)

	// Copy test file
	testFile := filepath.Join("testdata", "SavedVariables", "AddonProfilesDB.lua")
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	destFile := filepath.Join(savedVarsDir, "AddonProfilesDB.lua")
	if err := os.WriteFile(destFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	mgr := NewManager(tmpDir, account, 5)
	db, err := mgr.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() error = %v", err)
	}

	if db == nil {
		t.Fatal("LoadProfiles() returned nil database")
	}

	if len(db.Global.Profiles) == 0 {
		t.Error("Expected global profiles, got none")
	}
}

func TestLoadProfilesNoFile(t *testing.T) {
	// Create test structure without SavedVariables file
	tmpDir, err := os.MkdirTemp("", "wow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	account := "TestAccount"
	os.MkdirAll(filepath.Join(tmpDir, "WTF", "Account", account), 0755)

	mgr := NewManager(tmpDir, account, 5)
	db, err := mgr.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() error = %v", err)
	}

	if db == nil {
		t.Fatal("LoadProfiles() returned nil database")
	}

	// Should return empty database
	if len(db.Global.Profiles) != 0 {
		t.Error("Expected empty profiles, got some")
	}
}

func TestApplyProfile(t *testing.T) {
	// Create test structure
	tmpDir, err := os.MkdirTemp("", "wow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	account := "TestAccount"
	accountDir := filepath.Join(tmpDir, "WTF", "Account", account)
	os.MkdirAll(accountDir, 0755)

	// Create initial AddOns.txt
	addonsPath := filepath.Join(accountDir, "AddOns.txt")
	initialAddons := map[string]bool{
		"Addon1": true,
		"Addon2": false,
	}
	writeAddOnsFile(addonsPath, initialAddons)

	mgr := NewManager(tmpDir, account, 5)

	// Apply new profile
	profile := &lua.Profile{
		Name:  "TestProfile",
		Scope: "account",
		Addons: map[string]bool{
			"Addon3": true,
			"Addon4": true,
			"Addon5": false,
		},
	}

	if err := mgr.ApplyProfile(profile); err != nil {
		t.Fatalf("ApplyProfile() error = %v", err)
	}

	// Verify AddOns.txt was updated
	newAddons, err := parseAddOnsFile(addonsPath)
	if err != nil {
		t.Fatalf("parseAddOnsFile() error = %v", err)
	}

	if len(newAddons) != len(profile.Addons) {
		t.Errorf("Expected %d addons, got %d", len(profile.Addons), len(newAddons))
	}

	for name, enabled := range profile.Addons {
		if newAddons[name] != enabled {
			t.Errorf("Addon %s: got %v, want %v", name, newAddons[name], enabled)
		}
	}

	// Verify backup was created
	entries, _ := os.ReadDir(accountDir)
	backupFound := false
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) != ".txt" {
			backupFound = true
			break
		}
	}
	if !backupFound {
		t.Error("Backup file not created")
	}
}

func TestCleanupBackups(t *testing.T) {
	// Create test structure
	tmpDir, err := os.MkdirTemp("", "wow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	account := "TestAccount"
	accountDir := filepath.Join(tmpDir, "WTF", "Account", account)
	os.MkdirAll(accountDir, 0755)

	addonsPath := filepath.Join(accountDir, "AddOns.txt")

	// Create 10 backup files
	for i := 0; i < 10; i++ {
		backupPath := fmt.Sprintf("%s.backup.2024010%d_120000", addonsPath, i)
		os.WriteFile(backupPath, []byte("test"), 0644)
	}

	mgr := NewManager(tmpDir, account, 3) // Keep only 3 backups
	if err := mgr.cleanupBackups(addonsPath); err != nil {
		t.Fatalf("cleanupBackups() error = %v", err)
	}

	// Count remaining backups
	entries, _ := os.ReadDir(accountDir)
	backupCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) != ".txt" {
			backupCount++
		}
	}

	if backupCount != 3 {
		t.Errorf("Expected 3 backups, got %d", backupCount)
	}
}

func TestGetActiveAddons(t *testing.T) {
	// Create test structure
	tmpDir, err := os.MkdirTemp("", "wow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	account := "TestAccount"
	accountDir := filepath.Join(tmpDir, "WTF", "Account", account)
	os.MkdirAll(accountDir, 0755)

	// Copy test AddOns.txt
	testFile := filepath.Join("testdata", "AddOns.txt")
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	destFile := filepath.Join(accountDir, "AddOns.txt")
	if err := os.WriteFile(destFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	mgr := NewManager(tmpDir, account, 5)
	addons, err := mgr.GetActiveAddons()
	if err != nil {
		t.Fatalf("GetActiveAddons() error = %v", err)
	}

	if len(addons) == 0 {
		t.Error("Expected addons, got none")
	}

	// Verify specific addons
	if !addons["Ace3"] {
		t.Error("Expected Ace3 to be enabled")
	}

	if addons["AddonProfiles"] {
		t.Error("Expected AddonProfiles to be disabled")
	}
}
