# ascii-art-web-dockerize

![CI](https://github.com/teovaira/ascii-art-web-dockerize/actions/workflows/ci.yml/badge.svg)

## Description

ASCII Art Generator — a CLI tool and web application written in Go that converts text to ASCII art using three banner styles (standard, shadow, thinkertoy), with optional ANSI 24-bit color support for the CLI. The web interface allows users to type text, choose a banner, and receive the rendered ASCII art directly in the browser.

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
- Docker support — containerized deployment with a single command

## Installation

### Prerequisites

- Go 1.22.2 or higher *(for building from source)*
- Docker *(for containerized deployment)*

### Run with Docker

```bash
# Build the image and start the container (port 8080)
./docker-build.sh
# or step by step:
docker image build -f Dockerfile -t ascii-art-web-docker .
docker container run --publish 8080:8080 --detach --name dockerize ascii-art-web-docker
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

```bash
# Stop and remove the container
./docker-build.sh stop

# Stop container and remove image
./docker-build.sh clean
```

### Build from source

```bash
# Clone the repository
git clone https://github.com/teovaira/ascii-art-web-dockerize.git
cd ascii-art-web-dockerize

# Build CLI binary
go build -o bin/ascii-art ./cmd/ascii-art
# or: make build

# Build web server binary
go build -o bin/ascii-art-web ./cmd/ascii-art-web
# or: make build-web
```

## Usage

### Web server

**Via Docker (recommended):**
```bash
./docker-build.sh          # build image + start container on port 8080
# or: make docker-build && make docker-run
```

**From source** (must run from repository root — server reads `templates/` and `static/` as relative paths):
```bash
go run ./cmd/ascii-art-web
# or: make run-web

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
go test ./...
# or: make test

# Run with coverage
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
# or: make coverage

# Run linters
golangci-lint run ./...
# or: make lint

# Format code
gofmt -w .
# or: make fmt

# Run web server (from project root)
go run ./cmd/ascii-art-web
# or: make run-web

# Run CLI with color mode
cd cmd/ascii-art && go run . --color=red "Hello World"
# or: make run-color
```

### Project Structure

```
ascii-art-web-dockerize/
├── .github/
│   └── workflows/
│       ├── ci.yml             # CI workflow (test, lint, build)
│       └── release.yml        # Release workflow (cross-platform binaries)
├── .gitignore                 # Git ignore rules
├── .golangci.yml              # Linter configuration
├── Dockerfile                 # Multi-stage Docker build
├── docker-build.sh            # Helper script: build image + run container
├── LICENSE                    # Project license
├── Makefile                   # Build automation (incl. docker-* targets)
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

### Docker Commands

```bash
make docker-build   # Build the Docker image (ascii-art-web-docker)
make docker-run     # Run the container (dockerize) on port 8080
make docker-stop    # Stop and remove the container
make docker-clean   # Stop container + remove image
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

## Implementation Details

### Algorithm

1. **Parsing** — each banner file (standard, shadow, thinkertoy) is a fixed-format text file where every printable ASCII character (32–126) occupies exactly 8 lines. `parser.LoadBanner()` reads the file and builds a `map[rune][]string` — a character map where each key is a character and the value is its 8-line ASCII art representation.

2. **Rendering** — `renderer.ASCII()` splits the input on newlines and processes each line word by word. For each line it iterates character by character, looks up the 8-line art block in the character map, and appends each row side by side using a `strings.Builder`. The result is the full multi-line ASCII art string.

3. **Web flow** — the browser sends a `POST /ascii-art` request with the form fields `text` and `banner`. The `HandleASCIIArt` handler validates the input via `internal/validation`, calls `GenerateASCII` which runs the parser and renderer, then re-renders `index.html` with the result embedded in a `<pre>` block. On error the same page is re-rendered with an error message and the appropriate HTTP status code (400 for bad input, 404 for unknown banner, 500 for render failure).

4. **Embedded banners** — banner files are compiled into the binary at build time using Go's `//go:embed` directive (`internal/banners`). Neither the web server nor the CLI require banner files on disk at runtime.

## Authors

| Name | Role |
|------|------|
| Theodore | Team Lead |
| Krystallenia | Backend |
| Vasiliki | Frontend |

## License

See [LICENSE](LICENSE) file for details.
