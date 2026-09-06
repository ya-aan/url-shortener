package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ya-aan/url-shortener/internal/http-server/handlers"
	"github.com/ya-aan/url-shortener/internal/service"
	"github.com/ya-aan/url-shortener/internal/storage"
	"github.com/ya-aan/url-shortener/internal/storage/postgres"
)

func TestPostgres_ConcurrentCreateURL(t *testing.T) {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("set TEST_POSTGRES_DSN to run the PostgreSQL integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	admin, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close(context.Background())

	schema := fmt.Sprintf("url_shortener_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _, err := admin.Exec(cleanupCtx, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Errorf("cleanup schema: %v", err)
		}
	}()

	if _, err := admin.Exec(ctx, "SET search_path TO "+schema); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../migrations/001_init.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := admin.Exec(ctx, string(migration)); err != nil {
		t.Fatal(err)
	}
	connectionURL, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	query := connectionURL.Query()
	query.Set("search_path", schema)
	connectionURL.RawQuery = query.Encode()
	db, err := postgres.New(ctx, connectionURL.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := handlers.Save(service.New(db))
	const requests = 16
	start := make(chan struct{})
	statuses := make(chan int, requests)
	for range requests {
		go func() {
			<-start
			req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader(`{"url":"https://example.com","alias":"shared"}`)).WithContext(ctx)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			statuses <- rec.Code
		}()
	}
	close(start)
	counts := make(map[int]int)
	for range requests {
		counts[<-statuses]++
	}
	if counts[http.StatusCreated] != 1 || counts[http.StatusConflict] != requests-1 {
		t.Fatalf("expected one 201 and %d conflicts, got %v", requests-1, counts)
	}
	if rawURL, err := db.GetURLByAlias(ctx, "shared"); err != nil || rawURL != "https://example.com" {
		t.Fatalf("saved URL: %q, error: %v", rawURL, err)
	}
	if _, err := db.SaveURL(ctx, "https://example.org", "shared"); !errors.Is(err, storage.ErrAliasExists) {
		t.Fatalf("expected storage.ErrAliasExists, got %v", err)
	}

	if _, err := admin.Exec(ctx, "SELECT setval('urls_id_seq', (SELECT id FROM urls WHERE alias = 'shared'), false)"); err != nil {
		t.Fatal(err)
	}
	_, err = db.SaveURL(ctx, "https://example.org", "different")
	var pgErr *pgconn.PgError
	if errors.Is(err, storage.ErrAliasExists) || !errors.As(err, &pgErr) || pgErr.ConstraintName != "urls_pkey" {
		t.Fatalf("unrelated unique violation must remain a database error, got %v", err)
	}

	db.Close()
	if _, err := db.GetURLByAlias(ctx, "shared"); err == nil {
		t.Fatal("expected the pool to be closed")
	}
}
