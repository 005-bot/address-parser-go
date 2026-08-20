# address-parser-go

Shared address parsing library: normalizes free-form street names (Russian) to
canonical street records with confidence scoring.

Extracted from [005-bot/monitor-go](https://github.com/005-bot/monitor-go)
(`internal/parser/address`) and published as a standalone Go module for reuse.

## Features

- Embedded SQLite street database (`streets.db`, loaded via pure-Go
  `modernc.org/sqlite`, no CGO required).
- Exact-match lookup with confidence `1.0`.
- Fuzzy matching: Levenshtein similarity + LCS blend (`0.3` / `0.7`),
  minimum confidence `0.6`, minimum LCS coverage of the stored name `0.4`.
- Input cleaning: lowercase, strip punctuation (Unicode `[^\p{L}\p{N}\s\-]`),
  collapse whitespace.
- `fx` module integration (`address.Module()`).

## API

- `type Config struct { DBPath string }` - optional path to an external
  streets DB (uses embedded copy when empty).
- `NewParser(cfg Config) (*Parser, error)` - load streets into memory.
- `(*Parser).Normalize(ctx context.Context, raw string) (*Match, error)` -
  match raw input to a street.
- `(*Parser).Stop()` - release temp resources (embedded DB extraction dir).
- `type Match struct { Name string; NormalizedName string; Confidence float64 }`.
- `var ErrNoMatch` - returned when no street scores at least `0.6`.

## Usage

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

As an fx module:

```go
import "github.com/005-bot/address-parser-go"

app := fx.New(
	address.Module(),
	fx.Provide(func() address.Config { return address.Config{} }),
)
```

## Development

- `make test` - race-enabled tests with coverage.
- `CGO_ENABLED=0 go build ./...` - pure-Go build (verified).
- `make lint` - golangci-lint.

## Attribution

- Parser code extracted from [005-bot/monitor-go](https://github.com/005-bot/monitor-go)
  (`internal/parser/address`), Apache-2.0.
- `streets.db` from [005-bot/address-parser](https://github.com/005-bot/address-parser),
  Apache-2.0; byte-identical to the Python address-parser database
  (MD5 `10072cee7eb84361125cbdaf76559093`).
- Fuzzy scoring (edlib Levenshtein + LCS blend) intentionally differs from the
  Python `difflib` scores; exact matches produce identical results.

## License

Apache-2.0 - see [LICENSE](LICENSE).
