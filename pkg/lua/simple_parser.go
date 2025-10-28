package lua

import (
	"fmt"
	"strconv"
)

// Simple recursive descent parser for WoW SavedVariables format
// This is specifically tailored to our AddonProfilesDB structure

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenLBrace
	tokenRBrace
	tokenLBracket
	tokenRBracket
	tokenComma
	tokenEquals
	tokenString
	tokenNumber
	tokenBool
	tokenNil
	tokenIdent
)

type token struct {
	typ   tokenType
	value string
}

type lexer struct {
	input  string
	pos    int
	tokens []token
}

func newLexer(input string) *lexer {
	return &lexer{
		input: input,
		pos:   0,
	}
}

func (l *lexer) lex() []token {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Skip whitespace and comments
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.pos++
			continue
		}

		if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
			// Skip comment
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			continue
		}

		switch ch {
		case '{':
			l.tokens = append(l.tokens, token{tokenLBrace, "{"})
			l.pos++
		case '}':
			l.tokens = append(l.tokens, token{tokenRBrace, "}"})
			l.pos++
		case '[':
			l.tokens = append(l.tokens, token{tokenLBracket, "["})
			l.pos++
		case ']':
			l.tokens = append(l.tokens, token{tokenRBracket, "]"})
			l.pos++
		case ',':
			l.tokens = append(l.tokens, token{tokenComma, ","})
			l.pos++
		case '=':
			l.tokens = append(l.tokens, token{tokenEquals, "="})
			l.pos++
		case '"':
			l.lexString()
		default:
			if isDigit(ch) || (ch == '-' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1])) {
				l.lexNumber()
			} else if isAlpha(ch) {
				l.lexIdent()
			} else {
				l.pos++
			}
		}
	}

	l.tokens = append(l.tokens, token{tokenEOF, ""})
	return l.tokens
}

func (l *lexer) lexString() {
	l.pos++ // skip opening quote
	start := l.pos

	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\\' {
			l.pos++ // skip escape
		}
		l.pos++
	}

	value := l.input[start:l.pos]
	l.pos++ // skip closing quote
	l.tokens = append(l.tokens, token{tokenString, value})
}

func (l *lexer) lexNumber() {
	start := l.pos

	if l.input[l.pos] == '-' {
		l.pos++
	}

	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}

	value := l.input[start:l.pos]
	l.tokens = append(l.tokens, token{tokenNumber, value})
}

func (l *lexer) lexIdent() {
	start := l.pos

	for l.pos < len(l.input) && (isAlpha(l.input[l.pos]) || isDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
		l.pos++
	}

	value := l.input[start:l.pos]

	// Check for keywords
	switch value {
	case "true":
		l.tokens = append(l.tokens, token{tokenBool, "true"})
	case "false":
		l.tokens = append(l.tokens, token{tokenBool, "false"})
	case "nil":
		l.tokens = append(l.tokens, token{tokenNil, "nil"})
	default:
		l.tokens = append(l.tokens, token{tokenIdent, value})
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

type parser struct {
	tokens []token
	pos    int
}

func newParser(tokens []token) *parser {
	return &parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *parser) peek() token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token{tokenEOF, ""}
}

func (p *parser) next() token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *parser) expect(typ tokenType) (token, error) {
	tok := p.next()
	if tok.typ != typ {
		return tok, fmt.Errorf("expected %v, got %v", typ, tok.typ)
	}
	return tok, nil
}

func (p *parser) parseTable() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if _, err := p.expect(tokenLBrace); err != nil {
		return nil, err
	}

	for p.peek().typ != tokenRBrace && p.peek().typ != tokenEOF {
		// Parse key
		if p.peek().typ == tokenLBracket {
			p.next() // consume [
			keyTok, err := p.expect(tokenString)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(tokenRBracket); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokenEquals); err != nil {
				return nil, err
			}

			// Parse value
			value, err := p.parseValue()
			if err != nil {
				return nil, err
			}

			result[keyTok.value] = value
		}

		// Optional comma
		if p.peek().typ == tokenComma {
			p.next()
		}
	}

	if _, err := p.expect(tokenRBrace); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *parser) parseValue() (interface{}, error) {
	tok := p.peek()

	switch tok.typ {
	case tokenLBrace:
		return p.parseTable()
	case tokenString:
		p.next()
		return tok.value, nil
	case tokenNumber:
		p.next()
		num, _ := strconv.ParseInt(tok.value, 10, 64)
		return num, nil
	case tokenBool:
		p.next()
		return tok.value == "true", nil
	case tokenNil:
		p.next()
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", tok.typ)
	}
}

// ParseSimple uses the simple parser
func ParseSimple(content string) (*Database, error) {
	// Tokenize
	lexer := newLexer(content)
	tokens := lexer.lex()

	// Parse
	parser := newParser(tokens)

	// Skip to AddonProfilesDB =
	for parser.peek().typ != tokenEOF {
		if parser.peek().typ == tokenIdent && parser.peek().value == "AddonProfilesDB" {
			parser.next() // consume identifier
			parser.expect(tokenEquals)
			break
		}
		parser.next()
	}

	// Parse the main table
	mainTable, err := parser.parseTable()
	if err != nil {
		return nil, err
	}

	// Convert to Database structure
	db := &Database{}
	db.Global.Profiles = make(map[string]*Profile)
	db.Char = make(map[string]struct {
		ActiveProfile string
		Profiles      map[string]*Profile
	})

	// Parse global section
	if globalRaw, ok := mainTable["global"]; ok {
		if globalMap, ok := globalRaw.(map[string]interface{}); ok {
			// Active profile
			if activeRaw, ok := globalMap["activeProfile"]; ok {
				if active, ok := activeRaw.(string); ok {
					db.Global.ActiveProfile = active
				}
			}

			// Profiles
			if profilesRaw, ok := globalMap["profiles"]; ok {
				if profilesMap, ok := profilesRaw.(map[string]interface{}); ok {
					for profileName, profileRaw := range profilesMap {
						if profileMap, ok := profileRaw.(map[string]interface{}); ok {
							profile := convertToProfile(profileName, "account", profileMap)
							db.Global.Profiles[profileName] = profile
						}
					}
				}
			}
		}
	}

	// Parse char section
	if charRaw, ok := mainTable["char"]; ok {
		if charMap, ok := charRaw.(map[string]interface{}); ok {
			for charKey, charDataRaw := range charMap {
				if charDataMap, ok := charDataRaw.(map[string]interface{}); ok {
					charData := struct {
						ActiveProfile string
						Profiles      map[string]*Profile
					}{
						Profiles: make(map[string]*Profile),
					}

					// Active profile
					if activeRaw, ok := charDataMap["activeProfile"]; ok {
						if active, ok := activeRaw.(string); ok {
							charData.ActiveProfile = active
						}
					}

					// Profiles
					if profilesRaw, ok := charDataMap["profiles"]; ok {
						if profilesMap, ok := profilesRaw.(map[string]interface{}); ok {
							for profileName, profileRaw := range profilesMap {
								if profileMap, ok := profileRaw.(map[string]interface{}); ok {
									profile := convertToProfile(profileName, "character", profileMap)
									charData.Profiles[profileName] = profile
								}
							}
						}
					}

					db.Char[charKey] = charData
				}
			}
		}
	}

	return db, nil
}

func convertToProfile(name, scope string, profileMap map[string]interface{}) *Profile {
	profile := &Profile{
		Name:     name,
		Scope:    scope,
		Addons:   make(map[string]bool),
		AutoDeps: true, // default
	}

	// Parse addons
	if addonsRaw, ok := profileMap["addons"]; ok {
		if addonsMap, ok := addonsRaw.(map[string]interface{}); ok {
			for addonName, enabledRaw := range addonsMap {
				if enabled, ok := enabledRaw.(bool); ok {
					profile.Addons[addonName] = enabled
				}
			}
		}
	}

	// Parse autoDeps
	if autoDepsRaw, ok := profileMap["autoDeps"]; ok {
		if autoDeps, ok := autoDepsRaw.(bool); ok {
			profile.AutoDeps = autoDeps
		}
	}

	// Parse created
	if createdRaw, ok := profileMap["created"]; ok {
		if created, ok := createdRaw.(int64); ok {
			profile.Created = created
		}
	}

	return profile
}
