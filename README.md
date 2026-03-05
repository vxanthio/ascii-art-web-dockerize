# ascii-art-web

![CI](https://github.com/teovaira/ascii-art-web/actions/workflows/ci.yml/badge.svg)

ASCII Art Generator — CLI tool and web interface for converting text to ASCII art using banner styles (standard, shadow, thinkertoy), with optional ANSI 24-bit color support.

## Features

- Web interface for browser-based ASCII art generation
- Three banner styles: standard, shadow, thinkertoy
- ANSI 24-bit color support (named colors, hex, RGB)
- Substring coloring for highlighting specific parts of the output
- User-friendly error feedback in the web UI
- High performance (sub-millisecond rendering)
- Zero external dependencies (Go standard library only)
- Cross-platform support (Linux, macOS, Windows)
- Support for newline characters in input

## Installation

### Prerequisites

- Go 1.22.2 or higher

### Build from source

```bash
# Clone the repository
git clone https://github.com/teovaira/ascii-art-web.git
cd ascii-art-web

# Build CLI binary
make build
# or: go build -o bin/ascii-art ./cmd/ascii-art

# Build web server binary
make build-web
# or: go build -o bin/ascii-art-web ./cmd/ascii-art-web
```

## Usage

### Web server

> **Note**: The web server binary and `go run` must be executed from the **repository root**. The server reads `templates/` and `static/` as relative paths at runtime. Unlike the CLI binary, it is not relocatable.

```bash
# Run from repository root
make run-web
# or: go run ./cmd/ascii-art-web

# Run the built binary — must also be from repository root
./bin/ascii-art-web

# Custom port
PORT=9090 go run ./cmd/ascii-art-web
```

Open [http://localhost:8080](http://localhost:8080) in your browser, type text, choose a banner, and submit.

### CLI — Normal mode

```bash
cd cmd/ascii-art && go run . "text" [banner]
```

### CLI — Color mode

```bash
cd cmd/ascii-art && go run . --color=<color> "text" [banner]
cd cmd/ascii-art && go run . --color=<color> <substring> "text" [banner]
```

**Arguments**:
- `text`: The text to convert to ASCII art (required)
- `banner`: Banner style - standard, shadow, or thinkertoy (optional, defaults to standard)
- `--color=<color>`: Color specification (optional)
- `substring`: Substring to colorize (optional, colors full text if omitted)

### Color formats

- **Named colors**: red, green, blue, yellow, cyan, magenta, white, black, orange, purple, pink, brown, gray
- **Hex**: `#RRGGBB` (e.g. `#ff0000`)
- **RGB**: `rgb(R,G,B)` (e.g. `rgb(255,0,0)`)

> **Note**: RGB format requires quoting or escaping in bash/zsh due to parentheses. Use single quotes (`'rgb(...)'`), double quotes (`"rgb(...)"`), or escape parentheses (`rgb\(...\)`).

### Examples

**Standard banner (default):**
```bash
cd cmd/ascii-art && go run . "Hello"
```

**Shadow banner:**
```bash
cd cmd/ascii-art && go run . "Hello" shadow
```

**Thinkertoy banner:**
```bash
cd cmd/ascii-art && go run . "Hello" thinkertoy
```

**Full text colored red:**
```bash
cd cmd/ascii-art && go run . --color=red "Hello World"
```

**Substring colored orange:**
```bash
cd cmd/ascii-art && go run . --color=orange GuYs "HeY GuYs"
```

**Single letter colored blue:**
```bash
cd cmd/ascii-art && go run . --color=blue B "RGB()"
```

**Hex color format:**
```bash
cd cmd/ascii-art && go run . --color=#ff0000 "Hello"
```

**RGB color format (with single quotes):**
```bash
cd cmd/ascii-art && go run . --color='rgb(255,0,0)' "Hello"
```

**Newline support:**
```bash
cd cmd/ascii-art && go run . "Hello\nWorld"
```

## Development

### Setup

```bash
# Run tests
make test

# Run with coverage
make coverage

# Run linters
make lint

# Format code
make fmt

# Run web server
make run-web

# Run CLI with color mode
make run-color
```

### Project Structure

```
ascii-art-web/
├── .github/
│   └── workflows/
│       ├── ci.yml             # CI workflow (test, lint, build)
│       └── release.yml        # Release workflow (cross-platform binaries)
├── .gitignore                 # Git ignore rules
├── .golangci.yml              # Linter configuration
├── LICENSE                    # Project license
├── Makefile                   # Build automation
├── go.mod                     # Go module file
├── AGENTS.md                  # AI agent instructions
├── CHANGELOG.md               # Version history
├── CONTRIBUTING.md            # Contribution guidelines
├── PERMISSIONS.md             # Team permissions
├── README.md                  # This file
├── diagrams/                  # Mermaid architecture diagrams
│   ├── architecture.md        # High-level system overview
│   ├── class-diagram.md       # Package types and relationships
│   ├── flowchart.md           # Program execution flow
│   └── sequence-diagram.md    # Color mode call sequence
├── cmd/
│   ├── ascii-art/             # CLI entry point
│   │   ├── main.go
│   │   ├── args.go
│   │   ├── banner.go
│   │   ├── color_mode.go
│   │   ├── main_test.go
│   │   ├── integration_test.go
│   │   └── testdata/          # Banner files and test fixtures
│   │       ├── standard.txt
│   │       ├── shadow.txt
│   │       ├── thinkertoy.txt
│   │       ├── corrupted.txt
│   │       ├── empty.txt
│   │       └── oversized.txt
│   └── ascii-art-web/         # Web server entry point
│       ├── main.go
│       └── integration_test.go
├── static/                    # Static web assets
│   ├── style.css
│   └── favicon files
├── templates/                 # HTML templates
│   ├── base.html
│   └── index.html
└── internal/
    ├── banners/               # Embedded banner files
    │   ├── banners.go
    │   ├── standard.txt
    │   ├── shadow.txt
    │   └── thinkertoy.txt
    ├── color/                 # Color specification parsing
    │   ├── color.go
    │   └── color_test.go
    ├── coloring/              # ANSI color application to ASCII art
    │   ├── coloring.go
    │   └── coloring_test.go
    ├── flagparser/            # CLI argument validation
    │   ├── flagparser.go
    │   └── flagparser_test.go
    ├── handlers/              # HTTP handlers and template cache
    │   ├── handlers.go
    │   ├── handlers_test.go
    │   └── template_cache.go
    ├── parser/                # Banner file parsing
    │   ├── banner_parser.go
    │   └── parser_test.go
    ├── renderer/              # ASCII art rendering
    │   ├── renderer.go
    │   └── renderer_test.go
    └── validation/            # Input validation for web handler
        ├── validation.go
        └── validation_test.go
```

### Running Tests

```bash
# All tests
make test

# With coverage report
make coverage
```

### Build Commands

```bash
# Build CLI for current platform
make build

# Build web server
make build-web

# Build CLI for all platforms
make build-all

# Build CLI for specific platforms
make build-linux    # Linux (amd64 and arm64)
make build-darwin   # macOS (amd64 and arm64)
make build-windows  # Windows (amd64)
```

## Architecture

The project serves both a CLI tool and a web server from a shared set of internal packages:

- **main** (`cmd/ascii-art`): CLI interface and orchestration
- **main** (`cmd/ascii-art-web`): HTTP server entry point
- **handlers** (`internal/handlers`): HTTP handlers, ASCII generation, template cache
- **banners** (`internal/banners`): Banner files embedded into the binary at compile time
- **parser** (`internal/parser`): Banner file reading and character map building
- **renderer** (`internal/renderer`): Text-to-ASCII-art conversion
- **validation** (`internal/validation`): Web input validation (text length, banner name)
- **color** (`internal/color`): Color specification parsing (named, hex, RGB)
- **coloring** (`internal/coloring`): ANSI color application to rendered ASCII art
- **flagparser** (`internal/flagparser`): CLI argument validation

For visual diagrams see the [diagrams/](diagrams/) folder:
[Architecture Overview](diagrams/architecture.md) | [Flowchart](diagrams/flowchart.md) | [Class Diagram](diagrams/class-diagram.md) | [Sequence Diagram](diagrams/sequence-diagram.md)

## Test Coverage

- **validation**: 100.0%
- **coloring**: 100.0%
- **flagparser**: 100.0%
- **handlers**: 89.2%
- **color**: 97.7%
- **parser**: 95.0%
- **renderer**: 97.1%
- **main (cli)**: 39.1% (os.Exit prevents in-process coverage; tested via integration)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## Documentation

- [AGENTS.md](AGENTS.md) - AI agent instructions
- [CHANGELOG.md](CHANGELOG.md) - Version history
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [PERMISSIONS.md](PERMISSIONS.md) - Team permissions
- [diagrams/](diagrams/) - Mermaid architecture diagrams

## License

See [LICENSE](LICENSE) file for details.
