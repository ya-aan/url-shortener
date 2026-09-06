package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ya-aan/url-shortener/internal/storage"
)

type saveStorage struct {
	URLStorage
	save func(context.Context, string, string) (int64, error)
}

func (s saveStorage) SaveURL(ctx context.Context, rawURL, alias string) (int64, error) {
	return s.save(ctx, rawURL, alias)
}

func TestCreateURL_StorageErrors(t *testing.T) {
	otherErr := errors.New("database unavailable")
	for _, tc := range []struct {
		name string
		err  error
		want error
	}{
		{"conflict", fmt.Errorf("insert: %w", storage.ErrAliasExists), ErrAliasExists},
		{"other error", otherErr, otherErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			svc := New(saveStorage{save: func(context.Context, string, string) (int64, error) {
				calls++
				return 0, tc.err
			}})
			id, alias, err := svc.CreateURL(context.Background(), "https://example.com", "example")
			if !errors.Is(err, tc.want) || id != 0 || alias != "" || calls != 1 {
				t.Fatalf("got (%d, %q, %v), calls=%d; want error %v and one insert", id, alias, err, calls, tc.want)
			}
		})
	}
}

func TestCreateURL_ValidAlias(t *testing.T) {
	for _, input := range []string{"a", "Example_123-abc", strings.Repeat("a", 64), ""} {
		t.Run(input, func(t *testing.T) {
			var savedAlias string
			svc := New(saveStorage{save: func(_ context.Context, rawURL, alias string) (int64, error) {
				if rawURL != "https://example.com" {
					t.Fatalf("unexpected URL: %q", rawURL)
				}
				savedAlias = alias
				return 42, nil
			}})
			id, alias, err := svc.CreateURL(context.Background(), "https://example.com", input)
			if err != nil || id != 42 || alias != savedAlias {
				t.Fatalf("unexpected result: (%d, %q, %v), saved alias %q", id, alias, err, savedAlias)
			}
			if input != "" && alias != input {
				t.Fatalf("alias changed: got %q, want %q", alias, input)
			}
			if input == "" && (len(alias) != 6 || !aliasPattern.MatchString(alias)) {
				t.Fatalf("invalid generated alias: %q", alias)
			}
		})
	}
}

func TestService_InvalidAlias(t *testing.T) {
	svc := New(nil)
	ctx := context.Background()
	operations := map[string]func(string) error{
		"create": func(alias string) error {
			_, _, err := svc.CreateURL(ctx, "https://example.com", alias)
			return err
		},
		"get": func(alias string) error {
			_, err := svc.GetURL(ctx, alias)
			return err
		},
		"delete": func(alias string) error { return svc.DeleteURL(ctx, alias) },
		"update": func(alias string) error { return svc.UpdateURL(ctx, alias, "https://example.com") },
	}
	for name, operation := range operations {
		t.Run(name, func(t *testing.T) {
			for _, alias := range []string{"", " ", "has space", "a/b", "a?b", "a#b", "a%b", ".", "..", "ссылка", "a\n", strings.Repeat("a", 65)} {
				if name == "create" && alias == "" {
					continue
				}
				t.Run(alias, func(t *testing.T) {
					if err := operation(alias); !errors.Is(err, ErrInvalidAlias) {
						t.Fatalf("got %v, want ErrInvalidAlias", err)
					}
				})
			}
		})
	}
}
