# address-parser-go

Shared Go library that normalizes free-form (Russian) street names to canonical
street records with confidence scoring. Extracted from
[005-bot/monitor-go](https://github.com/005-bot/monitor-go) and published as a
standalone module for reuse.

## Table of Contents

- [About The Project](#about-the-project)
- [Features](#features)
- [Getting Started](#getting-started)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgments](#acknowledgments)

## About The Project

`address-parser-go` resolves messy user-typed street names (e.g. `ул. Ленина`,
`пр. Мира`) to the canonical name stored in a street database, returning a
confidence score. It is intended for services that ingest free-form address
input and need a stable, normalized representation.

Key points:

- Pure-Go implementation — no CGO, no external services.
- Embedded SQLite street database; ships with the module, no separate download.
- Safe for concurrent use; suitable for long-running services and `fx`-based
  applications.

## Features

- Embedded SQLite street database (`streets.db`, loaded via pure-Go
  `modernc.org/sqlite`, no CGO required).
- Exact-match lookup with confidence `1.0`.
- Fuzzy matching: Levenshtein similarity + LCS blend (`0.3` / `0.7`),
  minimum confidence `0.6`, minimum LCS coverage of the input `0.4`.
- Input cleaning: lowercase, strip punctuation (Unicode `[^\p{L}\p{N}\s\-]`),
  collapse whitespace.
- `fx` module integration (`address.Module()`).
- Thread-safe `Parser` for concurrent normalization.

## Getting Started

### Prerequisites

- Go **1.25.7** or newer (see [`go.mod`](go.mod)).
- No external services, system libraries, or CGO toolchain required.

### Supported Environments

Any platform supported by the Go toolchain and `modernc.org/sqlite`
(pure-Go). Verified to build with `CGO_ENABLED=0`.

## Installation

Add the module to your project:

```sh
go get github.com/005-bot/address-parser-go
```

The embedded `streets.db` is included automatically — no additional asset
download is required.

## Usage

### Basic example

```go
package main

import (
	"context"
	"fmt"

	"github.com/005-bot/address-parser-go"
)

func main() {
	p, err := address.NewParser(address.Config{})
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	m, err := p.Normalize(context.Background(), "ул. Ленина")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s (%.2f)\n", m.Name, m.Confidence)
	// улица Ленина (1.00)
}
```

### As an fx module

```go
import "github.com/005-bot/address-parser-go"

app := fx.New(
	address.Module(),
	fx.Provide(func() address.Config { return address.Config{} }),
)
```

### API Reference

| Symbol                                                                         | Description                                                                                              |
| ------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------- |
| `type Config struct { DBPath string }`                                         | Optional path to an external streets DB (uses the embedded copy when empty). The koanf key is `db_path`. |
| `NewParser(cfg Config) (*Parser, error)`                                       | Load streets into memory and return a ready-to-use `Parser`.                                             |
| `(*Parser).Normalize(ctx, raw) (*Match, error)`                                | Match raw input to a street; returns `ErrNoMatch` when nothing scores at least `0.6`.                    |
| `(*Parser).Stop()`                                                             | Release temporary resources (the extracted embedded DB directory). Safe to call multiple times.          |
| `type Match struct { Name string; NormalizedName string; Confidence float64 }` | A resolved street with a confidence score in `[0, 1]`.                                                   |
| `var ErrNoMatch`                                                               | Returned by `Normalize` when no street reaches the minimum confidence threshold.                         |

## Configuration

The parser is configured via the `Config` struct passed to `NewParser`. There
are no environment variables; when used with [`koanf`](https://github.com/knadh/koanf),
the field is bound to the `db_path` key.

| Field    | Type     | Default         | Description                                                                                                                                |
| -------- | -------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `DBPath` | `string` | `""` (embedded) | Path to an external SQLite streets database. When empty, the embedded `streets.db` is extracted to a temporary directory and used instead. |

To use an external database:

```go
p, err := address.NewParser(address.Config{DBPath: "/path/to/streets.db"})
```

## Roadmap

Planned improvements and feature requests are tracked as
[GitHub issues](https://github.com/005-bot/address-parser-go/issues). See the
issue tracker for the current roadmap and to propose changes.

## Contributing

Contributions are welcome via pull requests against the `master` branch.

1. Fork the repository and create a feature branch.
2. Make your changes, following the existing code style.
3. Ensure `make lint` (golangci-lint) and `make test` (`go test -race` with
   coverage) pass locally.
4. Open a pull request against `master`.

Continuous integration runs on every pull request to `master`:
[golangci-lint](https://github.com/005-bot/address-parser-go/actions) for
linting and `go test -race` for tests. Inactive issues and PRs are
automatically closed after a period of inactivity.

## License

Distributed under the Apache-2.0 License. See [`LICENSE`](LICENSE) for details.

## Contact

- Repository: [github.com/005-bot/address-parser-go](https://github.com/005-bot/address-parser-go)
- Source (upstream): [005-bot/monitor-go](https://github.com/005-bot/monitor-go)
  (`internal/parser/address`)

## Acknowledgments

- Parser code extracted from
  [005-bot/monitor-go](https://github.com/005-bot/monitor-go)
  (`internal/parser/address`), Apache-2.0.
- `streets.db` from
  [005-bot/address-parser](https://github.com/005-bot/address-parser),
  Apache-2.0; byte-identical to the Python address-parser database
  (MD5 `10072cee7eb84361125cbdaf76559093`).
- Fuzzy scoring uses [hbollon/go-edlib](https://github.com/hbollon/go-edlib)
  (Levenshtein + LCS).
- Dependency injection via [uber-go/fx](https://github.com/uber-go/fx).
