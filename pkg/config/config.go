package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the application configuration
type Config struct {
	WowInstallPath  string `json:"wow_install_path"`
	SelectedAccount string `json:"selected_account"`
	BackupCount     int    `json:"backup_count"`
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		WowInstallPath:  "",
		SelectedAccount: "",
		BackupCount:     5,
	}
}

// GetConfigPath returns the OS-specific configuration file path
func GetConfigPath() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		configDir = filepath.Join(configDir, "AddonProfiles")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support", "AddonProfiles")
	default: // linux and others
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			configDir = filepath.Join(homeDir, ".config", "addonprofiles")
		} else {
			configDir = filepath.Join(configDir, "addonprofiles")
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing fields
	if config.BackupCount == 0 {
		config.BackupCount = 5
	}

	return &config, nil
}

// Save saves the configuration to the config file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.WowInstallPath == "" {
		return fmt.Errorf("WoW installation path is not set")
	}

	// Check if WoW directory exists
	if _, err := os.Stat(c.WowInstallPath); os.IsNotExist(err) {
		return fmt.Errorf("WoW installation path does not exist: %s", c.WowInstallPath)
	}

	// Check if WTF directory exists
	wtfPath := filepath.Join(c.WowInstallPath, "WTF")
	if _, err := os.Stat(wtfPath); os.IsNotExist(err) {
		return fmt.Errorf("WTF directory not found at: %s", wtfPath)
	}

	if c.BackupCount < 1 {
		return fmt.Errorf("backup count must be at least 1")
	}

	return nil
}
