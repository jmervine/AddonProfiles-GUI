package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.WowInstallPath != "" {
		t.Errorf("Expected empty WowInstallPath, got %s", cfg.WowInstallPath)
	}

	if cfg.SelectedAccount != "" {
		t.Errorf("Expected empty SelectedAccount, got %s", cfg.SelectedAccount)
	}

	if cfg.BackupCount != 5 {
		t.Errorf("Expected BackupCount 5, got %d", cfg.BackupCount)
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("GetConfigPath() returned empty path")
	}

	// Verify it's a valid path
	if !filepath.IsAbs(path) {
		t.Errorf("GetConfigPath() returned relative path: %s", path)
	}

	// Verify it ends with config.json
	if filepath.Base(path) != "config.json" {
		t.Errorf("Expected config.json, got %s", filepath.Base(path))
	}

	// OS-specific checks
	switch runtime.GOOS {
	case "darwin":
		if !containsPath(path, "Library/Application Support/AddonProfiles") {
			t.Errorf("macOS config path should contain Library/Application Support/AddonProfiles, got %s", path)
		}
	case "linux":
		if !containsPath(path, ".config/addonprofiles") && !containsPath(path, "/addonprofiles") {
			t.Errorf("Linux config path should contain .config/addonprofiles, got %s", path)
		}
	case "windows":
		if !containsPath(path, "AddonProfiles") {
			t.Errorf("Windows config path should contain AddonProfiles, got %s", path)
		}
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	// Remove config file if it exists
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}
	os.Remove(configPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Should return default config
	if cfg.BackupCount != 5 {
		t.Errorf("Expected default BackupCount 5, got %d", cfg.BackupCount)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create test config
	cfg := &Config{
		WowInstallPath:  "/path/to/wow",
		SelectedAccount: "TestAccount",
		BackupCount:     10,
	}

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Compare
	if loaded.WowInstallPath != cfg.WowInstallPath {
		t.Errorf("WowInstallPath = %v, want %v", loaded.WowInstallPath, cfg.WowInstallPath)
	}

	if loaded.SelectedAccount != cfg.SelectedAccount {
		t.Errorf("SelectedAccount = %v, want %v", loaded.SelectedAccount, cfg.SelectedAccount)
	}

	if loaded.BackupCount != cfg.BackupCount {
		t.Errorf("BackupCount = %v, want %v", loaded.BackupCount, cfg.BackupCount)
	}

	// Cleanup
	configPath, _ := GetConfigPath()
	os.Remove(configPath)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		setup   func() string
		cleanup func(string)
	}{
		{
			name: "empty WoW path",
			config: &Config{
				WowInstallPath:  "",
				SelectedAccount: "",
				BackupCount:     5,
			},
			wantErr: true,
		},
		{
			name: "non-existent path",
			config: &Config{
				WowInstallPath:  "/non/existent/path",
				SelectedAccount: "",
				BackupCount:     5,
			},
			wantErr: true,
		},
		{
			name: "path without WTF directory",
			config: &Config{
				WowInstallPath:  "",
				SelectedAccount: "",
				BackupCount:     5,
			},
			wantErr: true,
			setup: func() string {
				tmpDir, _ := os.MkdirTemp("", "wow-test-*")
				return tmpDir
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
		{
			name: "valid config",
			config: &Config{
				WowInstallPath:  "",
				SelectedAccount: "TestAccount",
				BackupCount:     5,
			},
			wantErr: false,
			setup: func() string {
				tmpDir, _ := os.MkdirTemp("", "wow-test-*")
				os.Mkdir(filepath.Join(tmpDir, "WTF"), 0755)
				return tmpDir
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
		{
			name: "invalid backup count",
			config: &Config{
				WowInstallPath:  "",
				SelectedAccount: "",
				BackupCount:     0,
			},
			wantErr: true,
			setup: func() string {
				tmpDir, _ := os.MkdirTemp("", "wow-test-*")
				os.Mkdir(filepath.Join(tmpDir, "WTF"), 0755)
				return tmpDir
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				dir := tt.setup()
				tt.config.WowInstallPath = dir
				defer tt.cleanup(dir)
			}

			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadWithMissingBackupCount(t *testing.T) {
	// Create config without backup count
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}

	// Write incomplete config
	data := []byte(`{
  "wow_install_path": "/path/to/wow",
  "selected_account": "TestAccount"
}`)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load should apply default
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.BackupCount != 5 {
		t.Errorf("Expected default BackupCount 5, got %d", cfg.BackupCount)
	}

	// Cleanup
	os.Remove(configPath)
}

// Helper function
func containsPath(path, substr string) bool {
	return filepath.ToSlash(path) != "" && 
		(filepath.ToSlash(path) == substr || 
		 len(filepath.ToSlash(path)) > len(substr) && 
		 (filepath.ToSlash(path)[len(filepath.ToSlash(path))-len(substr):] == substr ||
		  containsSubstring(filepath.ToSlash(path), substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

