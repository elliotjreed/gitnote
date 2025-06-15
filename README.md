# GitNote

A cross-platform CLI application for managing markdown notes in a Git repository with automatic naming conventions and directory-based organization.

## Features

- **Organized note creation** with interactive category selection
- **Automatic file naming** using `yyyy-mm-dd note title.md` format
- **Table of contents generation** for easy navigation
- **Powerful search functionality** across titles and content
- **Git integration** for version control
- **Cross-platform compatibility** (Windows, macOS, Linux)

## Installation

### From Source

```bash
git clone <repository-url>
cd gitnote
make build
```

### Using Go Install

```bash
go install github.com/your-username/gitnote@latest
```

## Usage

### Create a New Note

```bash
gitnote new
```

This command will:
1. Show available categories (directories) for selection
2. Allow navigation through subdirectories
3. Provide option to create new categories
4. Prompt for note title
5. Create a markdown file with automatic date formatting
6. Display the relative path of the created file

### Generate Index

```bash
gitnote index
```

Scans the directory for markdown notes and generates/updates `readme.md` with a table of contents. The structure follows the directory hierarchy:

```
## work

### management

[2025-01-05 managing expectations](/work/management/2025-01-05 managing expectations.md)
```

### Search Notes

```bash
# Search by title only
gitnote search "managing"

# Search in both title and content
gitnote search --full "managing"
```

Returns a list of notes matching the search term.

### Commit Changes

```bash
gitnote commit
```

Automatically adds and commits new or modified files with descriptive commit messages listing the changes.

### Pull Updates

```bash
gitnote pull
```

Pulls changes from the remote repository with automatic merge conflict handling. Provides options to:
- Manually resolve conflicts
- Roll back changes if conflicts occur

## Project Structure

```
gitnote/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command setup
│   ├── new.go          # Note creation command
│   ├── index.go        # Index generation command
│   ├── search.go       # Search command
│   ├── commit.go       # Git commit command
│   └── pull.go         # Git pull command
├── internal/           # Internal packages
│   ├── note/           # Note management
│   ├── git/            # Git operations
│   └── index/          # Index generation
├── main.go             # Application entry point
├── go.mod              # Go module definition
├── Makefile            # Build and test commands
└── README.md           # Documentation
```

## Development

### Building

```bash
make build
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint
```

### Cross-platform Builds

```bash
make build-all
```

This creates binaries for:
- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [PromptUI](https://github.com/manifoldco/promptui) - Interactive prompts

## Testing

The project includes comprehensive unit tests with 100% code coverage:

- **Unit tests** for all internal packages
- **Integration tests** for CLI commands
- **Git repository simulation** for testing git operations
- **Temporary directory setup** for isolated testing

Run tests with:

```bash
go test -v ./...
```

## Architecture

The application follows clean architecture principles:

- **Separation of concerns** between CLI, business logic, and external dependencies
- **Dependency injection** for testability
- **Interface-based design** for modularity
- **Error handling** with proper error wrapping
- **Cross-platform compatibility** using Go's standard library

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.