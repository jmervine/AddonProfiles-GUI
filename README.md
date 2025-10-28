# Addon Profile Manager

A cross-platform GUI application for managing World of Warcraft addon profiles outside of the game.

## Features

- **Profile Management**: View and apply addon profiles created with the in-game AddonProfiles addon
- **Safe Operations**: Automatic backups before modifying AddOns.txt
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Simple Interface**: Clean, easy-to-use GUI

## Installation

### From Source

```bash
go build -o bin/addonprofiles-manager ./cmd/gui
```

### Pre-built Binaries

Download the latest release for your platform from the releases page.

## Usage

1. Launch the application
2. Select your World of Warcraft installation directory when prompted
3. Browse your profiles in the left panel
4. Select a profile to view its addons in the middle panel
5. Click "Apply Profile" to activate the profile

## Building

### Prerequisites

- Go 1.19 or later
- Fyne dependencies (see [Fyne documentation](https://developer.fyne.io/started/))

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage
```

## Development

### Project Structure

```
├── cmd/gui/          # Main entry point
├── pkg/
│   ├── config/       # Configuration management
│   ├── lua/          # Lua SavedVariables parser
│   ├── ui/           # GUI components
│   └── wow/          # WoW data management
└── Makefile          # Build automation
```

### Testing

The project has comprehensive unit tests with >79% code coverage:

```bash
go test ./pkg/...
```

## Safety Features

- **Automatic Backups**: Creates timestamped backups before modifying AddOns.txt
- **Validation**: Verifies WoW directory structure before operations
- **Read-Only Profiles**: Only reads from SavedVariables, never writes
- **Confirmation Dialogs**: Confirms before applying profiles

## Related Projects

- [AddonProfiles](https://github.com/jmervine/AddonProfiles) - The in-game WoW addon

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please open an issue or pull request.

