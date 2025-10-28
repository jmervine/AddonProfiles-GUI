package lua

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Profile represents an addon profile
type Profile struct {
	Name     string
	Scope    string // "account" or "character"
	Addons   map[string]bool
	AutoDeps bool
	Created  int64
}

// Database represents the parsed AddonProfilesDB structure
type Database struct {
	Global struct {
		ActiveProfile string
		Profiles      map[string]*Profile
		Settings      map[string]interface{}
	}
	Char map[string]struct {
		ActiveProfile string
		Profiles      map[string]*Profile
	}
}

// Parser handles parsing of Lua SavedVariables files
type Parser struct {
	content string
	pos     int
}

// NewParser creates a new Lua parser
func NewParser(content string) *Parser {
	return &Parser{
		content: content,
		pos:     0,
	}
}

// ParseFile parses a Lua SavedVariables file
func ParseFile(filepath string) (*Database, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return Parse(string(content))
}

// Parse parses Lua content and extracts the database structure
func Parse(content string) (*Database, error) {
	// Try simple parser first
	db, err := ParseSimple(content)
	if err == nil {
		return db, nil
	}

	// Fallback to regex-based parser (kept for compatibility)
	return parseRegex(content)
}

// parseRegex is the old regex-based parser (kept as fallback)
func parseRegex(content string) (*Database, error) {
	db := &Database{}
	db.Global.Profiles = make(map[string]*Profile)
	db.Char = make(map[string]struct {
		ActiveProfile string
		Profiles      map[string]*Profile
	})

	// Parse global profiles
	globalProfiles, err := extractTable(content, "AddonProfilesDB", "global", "profiles")
	if err == nil {
		for profileName, profileData := range globalProfiles {
			profile, err := parseProfile(profileName, "account", profileData)
			if err == nil {
				db.Global.Profiles[profileName] = profile
			}
		}
	}

	// Parse global active profile
	if activeProfile := extractString(content, "AddonProfilesDB", "global", "activeProfile"); activeProfile != "" {
		db.Global.ActiveProfile = activeProfile
	}

	// Parse character profiles
	charSection := extractCharSection(content)
	for charKey, charContent := range charSection {
		charData := struct {
			ActiveProfile string
			Profiles      map[string]*Profile
		}{
			Profiles: make(map[string]*Profile),
		}

		charProfiles, err := extractTableFromContent(charContent, "profiles")
		if err == nil {
			for profileName, profileData := range charProfiles {
				profile, err := parseProfile(profileName, "character", profileData)
				if err == nil {
					charData.Profiles[profileName] = profile
				}
			}
		}

		if activeProfile := extractStringFromContent(charContent, "activeProfile"); activeProfile != "" {
			charData.ActiveProfile = activeProfile
		}

		db.Char[charKey] = charData
	}

	return db, nil
}

// parseProfile parses a profile table
func parseProfile(name, scope, content string) (*Profile, error) {
	profile := &Profile{
		Name:   name,
		Scope:  scope,
		Addons: make(map[string]bool),
	}

	// Parse addons table
	addonsTable, err := extractTableFromContent(content, "addons")
	if err == nil {
		for addonName, value := range addonsTable {
			if strings.Contains(value, "true") {
				profile.Addons[addonName] = true
			}
		}
	}

	// Parse autoDeps
	if strings.Contains(content, "autoDeps") {
		autoDepsStr := extractStringFromContent(content, "autoDeps")
		profile.AutoDeps = autoDepsStr == "true"
	} else {
		profile.AutoDeps = true // default
	}

	// Parse created timestamp
	if createdStr := extractStringFromContent(content, "created"); createdStr != "" {
		if created, err := strconv.ParseInt(createdStr, 10, 64); err == nil {
			profile.Created = created
		}
	}

	return profile, nil
}

// extractTable extracts a nested table from Lua content
func extractTable(content string, keys ...string) (map[string]string, error) {
	pattern := buildTablePattern(keys...)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return nil, fmt.Errorf("table not found: %v", keys)
	}

	return parseTableContent(matches[1]), nil
}

// extractTableFromContent extracts a table from already isolated content
func extractTableFromContent(content, key string) (map[string]string, error) {
	// Find the key and its table
	pattern := fmt.Sprintf(`\["%s"\]\s*=\s*\{([^}]*(?:\{[^}]*\}[^}]*)*)\}`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return nil, fmt.Errorf("table not found: %s", key)
	}

	return parseTableContent(matches[1]), nil
}

// parseTableContent parses the content of a Lua table
func parseTableContent(content string) map[string]string {
	result := make(map[string]string)

	// Match table entries: ["key"] = value or ["key"] = { ... }
	entryPattern := regexp.MustCompile(`\["([^"]+)"\]\s*=\s*(\{[^}]*(?:\{[^}]*\}[^}]*)*\}|[^,\n]+)`)
	matches := entryPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			key := match[1]
			value := strings.TrimSpace(match[2])
			value = strings.Trim(value, ",")
			result[key] = value
		}
	}

	return result
}

// extractString extracts a string value from Lua content
func extractString(content string, keys ...string) string {
	pattern := buildStringPattern(keys...)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

// extractStringFromContent extracts a string from already isolated content
func extractStringFromContent(content, key string) string {
	pattern := fmt.Sprintf(`\["%s"\]\s*=\s*"([^"]*)"`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)

	if len(matches) >= 2 {
		return matches[1]
	}

	// Try without quotes (for booleans/numbers)
	pattern = fmt.Sprintf(`\["%s"\]\s*=\s*([^,\n]+)`, regexp.QuoteMeta(key))
	re = regexp.MustCompile(pattern)
	matches = re.FindStringSubmatch(content)

	if len(matches) >= 2 {
		return strings.TrimSpace(strings.Trim(matches[1], ","))
	}

	return ""
}

// extractCharSection extracts all character sections
func extractCharSection(content string) map[string]string {
	result := make(map[string]string)

	// Find the char section
	charPattern := regexp.MustCompile(`\["char"\]\s*=\s*\{(.*)\}[\s\n]*\}[\s\n]*$`)
	charMatches := charPattern.FindStringSubmatch(content)

	if len(charMatches) < 2 {
		return result
	}

	charContent := charMatches[1]

	// Extract each character entry
	scanner := bufio.NewScanner(strings.NewReader(charContent))
	var currentChar string
	var currentContent strings.Builder
	bracketDepth := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Check for character key
		charKeyPattern := regexp.MustCompile(`\["([^"]+\s+-\s+[^"]+)"\]\s*=\s*\{`)
		if matches := charKeyPattern.FindStringSubmatch(line); len(matches) >= 2 {
			if currentChar != "" && currentContent.Len() > 0 {
				result[currentChar] = currentContent.String()
			}
			currentChar = matches[1]
			currentContent.Reset()
			bracketDepth = 1
			continue
		}

		if currentChar != "" {
			currentContent.WriteString(line)
			currentContent.WriteString("\n")

			// Track bracket depth
			bracketDepth += strings.Count(line, "{") - strings.Count(line, "}")

			if bracketDepth == 0 {
				result[currentChar] = currentContent.String()
				currentChar = ""
				currentContent.Reset()
			}
		}
	}

	if currentChar != "" && currentContent.Len() > 0 {
		result[currentChar] = currentContent.String()
	}

	return result
}

// buildTablePattern builds a regex pattern for nested table access
func buildTablePattern(keys ...string) string {
	pattern := regexp.QuoteMeta(keys[0])
	for i := 1; i < len(keys); i++ {
		pattern += fmt.Sprintf(`\s*=\s*\{[^}]*\["%s"\]`, regexp.QuoteMeta(keys[i]))
	}
	pattern += `\s*=\s*\{([^}]*(?:\{[^}]*\}[^}]*)*)\}`
	return pattern
}

// buildStringPattern builds a regex pattern for nested string access
func buildStringPattern(keys ...string) string {
	pattern := regexp.QuoteMeta(keys[0])
	for i := 1; i < len(keys)-1; i++ {
		pattern += fmt.Sprintf(`\s*=\s*\{[^}]*\["%s"\]`, regexp.QuoteMeta(keys[i]))
	}
	lastKey := keys[len(keys)-1]
	pattern += fmt.Sprintf(`\s*=\s*\{[^}]*\["%s"\]\s*=\s*"([^"]*)"`, regexp.QuoteMeta(lastKey))
	return pattern
}
