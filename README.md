# Addon Profile Manager

A cross-platform GUI application for managing World of Warcraft addon profiles outside of the game.

## Features

- **Profile Management**: View and apply addon profiles created with the in-game AddonProfiles addon
- **Safe Operations**: Automatic backups before modifying AddOns.txt
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Simple Interface**: Clean, easy-to-use GUI

## Installation

### Pre-built Binaries (Recommended)

Download the latest release for your platform:

👉 **[Download from Releases](https://github.com/jmervine/AddonProfiles-GUI/releases/latest)**

#### Windows
1. Download `addonprofiles-manager-windows-amd64.zip`
2. Extract the ZIP file
3. Run `addonprofiles-manager.exe`

**Note for Windows users:** Windows Defender may flag the executable as a potential threat (false positive). This is common with unsigned Go applications. To run the application:
- Click "More info" → "Run anyway" when Windows SmartScreen appears
- Or add an exception in Windows Defender: Settings → Virus & threat protection → Manage settings → Add or remove exclusions

#### macOS
1. Download: `addonprofiles-manager-macos-universal.tar.gz` (works on both Intel and Apple Silicon)
2. Extract the archive: `tar xzf addonprofiles-manager-macos-universal.tar.gz`
3. Make executable (if needed): `chmod +x addonprofiles-manager`
4. Run: `./addonprofiles-manager`
5. On first launch, you may need to right-click → Open to bypass Gatekeeper

#### Linux
1. Download `addonprofiles-manager-linux-amd64.tar.gz`
2. Extract: `tar xzf addonprofiles-manager-linux-amd64.tar.gz`
3. Make executable: `chmod +x addonprofiles-manager`
4. Run: `./addonprofiles-manager`

### From Source

```bash
git clone https://github.com/jmervine/AddonProfiles-GUI.git
cd AddonProfiles-GUI
make build
./bin/addonprofiles-manager
```

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
- **For cross-compilation**: Docker and `fyne-cross` tool

### Build Commands

```bash
# Build for current platform (native build)
make build

# Install fyne-cross for cross-compilation
make install-fyne-cross

# Build for specific platforms (requires Docker)
make build-windows  # Windows AMD64
make build-mac      # macOS Intel + ARM64
make build-linux    # Linux AMD64

# Build for all platforms at once (requires Docker)
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Cross-Compilation Notes

Fyne apps require `fyne-cross` for cross-platform builds because they use CGO:

1. **Install Docker**: Required by fyne-cross for cross-compilation
2. **Install fyne-cross**: `make install-fyne-cross` or `go install github.com/fyne-io/fyne-cross@latest`
3. **Build**: `make build-all` will create binaries in `fyne-cross/dist/` directory

**Alternative**: Build natively on each platform using `make build` if you don't have Docker.

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

