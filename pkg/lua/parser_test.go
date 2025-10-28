package lua

import (
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid profile file",
			file:    "valid_profile.lua",
			wantErr: false,
		},
		{
			name:    "empty profile file",
			file:    "empty_profile.lua",
			wantErr: false,
		},
		{
			name:    "non-existent file",
			file:    "does_not_exist.lua",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("testdata", tt.file)
			db, err := ParseFile(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && db == nil {
				t.Error("ParseFile() returned nil database without error")
			}
		})
	}
}

func TestParseValidProfile(t *testing.T) {
	path := filepath.Join("testdata", "valid_profile.lua")
	db, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	// Test global profiles
	if len(db.Global.Profiles) != 2 {
		t.Errorf("Expected 2 global profiles, got %d", len(db.Global.Profiles))
	}

	// Test Default profile
	defaultProfile, ok := db.Global.Profiles["Default"]
	if !ok {
		t.Fatal("Default profile not found")
	}

	if defaultProfile.Name != "Default" {
		t.Errorf("Expected profile name 'Default', got '%s'", defaultProfile.Name)
	}

	if defaultProfile.Scope != "account" {
		t.Errorf("Expected scope 'account', got '%s'", defaultProfile.Scope)
	}

	if !defaultProfile.AutoDeps {
		t.Error("Expected AutoDeps to be true")
	}

	if len(defaultProfile.Addons) != 3 {
		t.Errorf("Expected 3 addons in Default profile, got %d", len(defaultProfile.Addons))
	}

	expectedAddons := []string{"Ace3", "DBM-Core", "Details"}
	for _, addon := range expectedAddons {
		if !defaultProfile.Addons[addon] {
			t.Errorf("Expected addon '%s' to be enabled", addon)
		}
	}

	// Test Raiding profile
	raidingProfile, ok := db.Global.Profiles["Raiding"]
	if !ok {
		t.Fatal("Raiding profile not found")
	}

	if len(raidingProfile.Addons) != 3 {
		t.Errorf("Expected 3 addons in Raiding profile, got %d", len(raidingProfile.Addons))
	}

	// Test active profile
	if db.Global.ActiveProfile != "Default" {
		t.Errorf("Expected active profile 'Default', got '%s'", db.Global.ActiveProfile)
	}
}

func TestParseCharacterProfiles(t *testing.T) {
	path := filepath.Join("testdata", "valid_profile.lua")
	db, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	charKey := "TestChar - TestRealm"
	charData, ok := db.Char[charKey]
	if !ok {
		t.Fatalf("Character '%s' not found", charKey)
	}

	if len(charData.Profiles) != 1 {
		t.Errorf("Expected 1 character profile, got %d", len(charData.Profiles))
	}

	pvpProfile, ok := charData.Profiles["PvP"]
	if !ok {
		t.Fatal("PvP profile not found")
	}

	if pvpProfile.Scope != "character" {
		t.Errorf("Expected scope 'character', got '%s'", pvpProfile.Scope)
	}

	if pvpProfile.AutoDeps {
		t.Error("Expected AutoDeps to be false")
	}

	if len(pvpProfile.Addons) != 2 {
		t.Errorf("Expected 2 addons in PvP profile, got %d", len(pvpProfile.Addons))
	}

	if charData.ActiveProfile != "PvP" {
		t.Errorf("Expected active profile 'PvP', got '%s'", charData.ActiveProfile)
	}
}

func TestParseEmptyProfile(t *testing.T) {
	path := filepath.Join("testdata", "empty_profile.lua")
	db, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if len(db.Global.Profiles) != 0 {
		t.Errorf("Expected 0 global profiles, got %d", len(db.Global.Profiles))
	}

	if len(db.Char) != 0 {
		t.Errorf("Expected 0 character sections, got %d", len(db.Char))
	}

	if db.Global.ActiveProfile != "" {
		t.Errorf("Expected no active profile, got '%s'", db.Global.ActiveProfile)
	}
}

func TestParseMalformed(t *testing.T) {
	path := filepath.Join("testdata", "malformed.lua")
	db, err := ParseFile(path)

	// Should not error, but might return incomplete data
	if err != nil {
		t.Logf("ParseFile() returned error (expected for malformed): %v", err)
	}

	if db != nil {
		// If it parsed something, profiles should be incomplete
		t.Logf("Parsed malformed file, got %d global profiles", len(db.Global.Profiles))
	}
}

func TestParseProfile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantName  string
		wantScope string
		wantCount int
		wantDeps  bool
	}{
		{
			name: "simple profile",
			content: `
				["addons"] = {
					["Addon1"] = true,
					["Addon2"] = true,
				},
				["autoDeps"] = true,
			`,
			wantName:  "Test",
			wantScope: "account",
			wantCount: 2,
			wantDeps:  true,
		},
		{
			name: "profile with autoDeps false",
			content: `
				["addons"] = {
					["Addon1"] = true,
				},
				["autoDeps"] = false,
			`,
			wantName:  "Test",
			wantScope: "character",
			wantCount: 1,
			wantDeps:  false,
		},
		{
			name: "profile without autoDeps (should default to true)",
			content: `
				["addons"] = {
					["Addon1"] = true,
				},
			`,
			wantName:  "Test",
			wantScope: "account",
			wantCount: 1,
			wantDeps:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := parseProfile(tt.wantName, tt.wantScope, tt.content)
			if err != nil {
				t.Fatalf("parseProfile() error = %v", err)
			}

			if profile.Name != tt.wantName {
				t.Errorf("Name = %v, want %v", profile.Name, tt.wantName)
			}

			if profile.Scope != tt.wantScope {
				t.Errorf("Scope = %v, want %v", profile.Scope, tt.wantScope)
			}

			if len(profile.Addons) != tt.wantCount {
				t.Errorf("Addon count = %v, want %v", len(profile.Addons), tt.wantCount)
			}

			if profile.AutoDeps != tt.wantDeps {
				t.Errorf("AutoDeps = %v, want %v", profile.AutoDeps, tt.wantDeps)
			}
		})
	}
}

func TestExtractTableFromContent(t *testing.T) {
	content := `
		["addons"] = {
			["Addon1"] = true,
			["Addon2"] = true,
			["Addon3"] = true,
		},
		["autoDeps"] = true,
	`

	table, err := extractTableFromContent(content, "addons")
	if err != nil {
		t.Fatalf("extractTableFromContent() error = %v", err)
	}

	if len(table) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(table))
	}

	for i := 1; i <= 3; i++ {
		key := "Addon" + string(rune('0'+i))
		if _, ok := table[key]; !ok {
			t.Errorf("Expected key '%s' not found", key)
		}
	}
}

func TestExtractStringFromContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		key     string
		want    string
	}{
		{
			name:    "string value",
			content: `["activeProfile"] = "Default",`,
			key:     "activeProfile",
			want:    "Default",
		},
		{
			name:    "boolean true",
			content: `["autoDeps"] = true,`,
			key:     "autoDeps",
			want:    "true",
		},
		{
			name:    "boolean false",
			content: `["autoDeps"] = false,`,
			key:     "autoDeps",
			want:    "false",
		},
		{
			name:    "number",
			content: `["created"] = 1698765432,`,
			key:     "created",
			want:    "1698765432",
		},
		{
			name:    "missing key",
			content: `["other"] = "value",`,
			key:     "missing",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractStringFromContent(tt.content, tt.key)
			if got != tt.want {
				t.Errorf("extractStringFromContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTableContent(t *testing.T) {
	content := `
		["key1"] = "value1",
		["key2"] = true,
		["key3"] = {
			["nested"] = "data",
		},
	`

	result := parseTableContent(content)

	if len(result) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(result))
	}

	if result["key1"] != `"value1"` {
		t.Errorf("key1 = %v, want %v", result["key1"], `"value1"`)
	}

	if result["key2"] != "true" {
		t.Errorf("key2 = %v, want true", result["key2"])
	}
}
