# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
```bash
make build          # Build the gitnote binary to bin/gitnote
make run            # Run the application directly with go run
make install        # Install the binary system-wide
```

### Testing
```bash
make test           # Run all tests with verbose output
make test-coverage  # Generate test coverage report (coverage.html)
go test ./internal/note     # Run tests for specific package
go test -run TestName       # Run specific test
```

### Code Quality
```bash
make fmt            # Format all Go code
make lint           # Run golangci-lint (requires installation)
make deps           # Download and organize dependencies
```

### Cross-platform Builds
```bash
make build-all      # Build for Linux, macOS (amd64/arm64), and Windows
```

## Architecture Overview

GitNote follows a clean architecture pattern with clear separation between CLI interface, business logic, and external dependencies.

### Core Components

**CLI Layer (`cmd/`)**
- Uses Cobra framework for command structure
- Each command is a separate file (new.go, index.go, search.go, commit.go, pull.go)
- Commands delegate business logic to internal packages
- Interactive prompts handled with PromptUI

**Business Logic (`internal/`)**
- `note.Manager`: Core note operations (create, find, search)
- `git.Manager`: Git operations abstraction
- `index.Generator`: README generation logic
- All managers accept working directory for testability

**Key Patterns**
- Manager structs with dependency injection via constructors
- Error wrapping with context (`fmt.Errorf("context: %w", err)`)
- File operations use `filepath.Join()` for cross-platform paths
- Sorting applied consistently (categories, notes, search results)

### Note File Convention
- Format: `yyyy-mm-dd note title.md`
- Content starts with `# note title` heading
- Directory structure reflects categories/subcategories
- Hidden directories (starting with `.`) are ignored

### Testing Strategy
- Comprehensive unit tests with temporary directories
- Git operations tested with real git repositories in temp dirs
- CLI commands tested through direct function calls
- Test helpers for git repo setup and teardown

### Important Implementation Details
- Category navigation supports unlimited nesting depth
- Search can be title-only or full-content with `--full` flag
- Index generation creates hierarchical markdown with proper heading levels
- Git commit messages are auto-generated based on file changes
- Merge conflict handling provides user choice between manual resolution and rollback