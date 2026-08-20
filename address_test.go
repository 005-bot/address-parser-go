package address_test

import (
	"context"
	"errors"
	"testing"

	"github.com/005-bot/address-parser-go"
)

func newTestParser(t *testing.T) *address.Parser {
	t.Helper()

	p, err := address.NewParser(address.Config{})
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}

	t.Cleanup(p.Stop)

	return p
}

func TestNormalize(t *testing.T) {
	p := newTestParser(t)
	ctx := context.Background()

	tests := []struct {
		input string
		want  string
	}{
		{"ул. Ленина", "улица Ленина"},
		{"пр. Мира", "проспект Мира"},
		{"Советская", "Советская улица"},
		{"Кольцевая", "Кольцевая улица"},
		{"Красноярский рабочий", "проспект имени газеты Красноярский Рабочий"},
		{"улица Рокоссовского", "улица Рокоссовского"},
		{"Металлургов", "проспект Металлургов"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			m, nerr := p.Normalize(ctx, tt.input)
			if nerr != nil {
				t.Fatalf("Normalize(%q) error: %v", tt.input, nerr)
			}
			if m.Name != tt.want {
				t.Errorf("Normalize(%q) = %q (conf=%.4f), want %q", tt.input, m.Name, m.Confidence, tt.want)
			}
			if m.Confidence < 0.6 {
				t.Errorf("Normalize(%q) confidence=%.4f < 0.6", tt.input, m.Confidence)
			}
		})
	}
}

func TestNormalizeExactMatch(t *testing.T) {
	p := newTestParser(t)
	ctx := context.Background()

	m, err := p.Normalize(ctx, "улица Рокоссовского")
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	if m.Name != "улица Рокоссовского" {
		t.Errorf("Name = %q, want %q", m.Name, "улица Рокоссовского")
	}
	if m.Confidence != 1.0 {
		t.Errorf("Confidence = %.4f, want 1.0", m.Confidence)
	}
	if m.NormalizedName != "улица рокоссовского" {
		t.Errorf("NormalizedName = %q, want %q", m.NormalizedName, "улица рокоссовского")
	}
}

func TestNormalizeBelowThreshold(t *testing.T) {
	p := newTestParser(t)
	ctx := context.Background()

	tests := []string{
		"qqqqq",
		"вжпщрвукп",
		"!!!",
		"",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := p.Normalize(ctx, input)
			if err == nil {
				t.Fatalf("Normalize(%q) succeeded, want ErrNoMatch", input)
			}
			if !errors.Is(err, address.ErrNoMatch) {
				t.Errorf("Normalize(%q) error = %v, want ErrNoMatch", input, err)
			}
		})
	}
}

func TestNormalizeCaseInsensitiveExact(t *testing.T) {
	p := newTestParser(t)
	ctx := context.Background()

	m, err := p.Normalize(ctx, "УЛИЦА РОКОССОВСКОГО")
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	if m.Name != "улица Рокоссовского" {
		t.Errorf("Name = %q, want %q", m.Name, "улица Рокоссовского")
	}
}

func TestNormalizeCanceledContext(t *testing.T) {
	p := newTestParser(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := p.Normalize(ctx, "улица Рокоссовского"); !errors.Is(err, context.Canceled) {
		t.Errorf("Normalize(exact) error = %v, want context.Canceled", err)
	}

	if _, err := p.Normalize(ctx, "qqqqq"); !errors.Is(err, context.Canceled) {
		t.Errorf("Normalize(fuzzy) error = %v, want context.Canceled", err)
	}
}

type cancelAfterNContext struct {
	context.Context

	remaining int
	cancel    context.CancelFunc
}

func newCancelAfterNContext(n int) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &cancelAfterNContext{Context: ctx, remaining: n, cancel: cancel}
}

func (c *cancelAfterNContext) Err() error {
	if c.remaining <= 0 {
		c.cancel()
	} else {
		c.remaining--
	}
	return c.Context.Err()
}

func TestNormalizeCancelDuringFuzzyScan(t *testing.T) {
	p := newTestParser(t)

	ctx := newCancelAfterNContext(3)

	_, err := p.Normalize(ctx, "qqqqq")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Normalize error = %v, want context.Canceled", err)
	}
}

func TestNormalizeSingleRuneNoMatch(t *testing.T) {
	p := newTestParser(t)
	ctx := context.Background()

	for _, input := range []string{"a", "у"} {
		_, err := p.Normalize(ctx, input)
		if !errors.Is(err, address.ErrNoMatch) {
			t.Errorf("Normalize(%q) error = %v, want ErrNoMatch", input, err)
		}
	}
}
