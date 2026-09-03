package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ThanhNV121097/project-f5d06cd0/backend/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type server struct {
	db *pgxpool.Pool
}

type apiErrorDetail struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiError struct {
	Error struct {
		Code      string           `json:"code"`
		Message   string           `json:"message"`
		Details   []apiErrorDetail `json:"details,omitempty"`
		RequestID string           `json:"request_id"`
	} `json:"error"`
}


type greeting struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type helloResponse struct {
	Message string `json:"message"`
}

type greetingsResponse struct {
	Greetings  []greeting `json:"greetings"`
	NextCursor *string    `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	if err := applyMigrations(ctx, db); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	srv := server{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", srv.healthz)
	mux.HandleFunc("GET /v1/health", srv.apiHealth)
	mux.HandleFunc("GET /v1/hello", srv.apiHello)
	mux.HandleFunc("POST /v1/greetings", srv.createGreeting)
	mux.HandleFunc("GET /v1/greetings", srv.listGreetings)

	addr := ":" + listenPort()
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func listenPort() string { if port := os.Getenv("PORT"); port != "" { return port }; if port := os.Getenv("APP_PORT"); port != "" { return port }; return "8080" }

func (s server) healthz(w http.ResponseWriter, r *http.Request) { ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel(); if err := s.db.Ping(ctx); err != nil { http.Error(w, "database unavailable", http.StatusServiceUnavailable); return }; writeJSON(w, http.StatusOK, map[string]string{"status":"ok"}) }

func (s server) apiHealth(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, map[string]string{"status":"ok"}) }

func (s server) apiHello(w http.ResponseWriter, r *http.Request) {
	requestID := apiRequestID(r)
	addRequestID(w, requestID)
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if utf8.RuneCountInString(name) > 80 {
		writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed.", []apiErrorDetail{{Field: "name", Code: "TOO_LONG", Message: "Name must be 80 characters or fewer."}}, requestID)
		return
	}
	if name == "" {
		name = "World"
	}
	writeJSON(w, http.StatusOK, helloResponse{Message: "Hello, " + name + "!"})
}

func (s server) createGreeting(w http.ResponseWriter, r *http.Request) {
	requestID := apiRequestID(r)
	addRequestID(w, requestID)
	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.HasPrefix(ct, "application/json") { writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Content type must be application/json.", nil, requestID); return }
	body, err := io.ReadAll(io.LimitReader(r.Body, 16<<10+1))
	if err != nil { writeAPIError(w, http.StatusInternalServerError, "INTERNAL", "Unexpected error.", nil, requestID); return }
	if len(body) > 16<<10 { writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Body too large.", nil, requestID); return }
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil { writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Malformed JSON.", nil, requestID); return }
	name, ok := parseStringField(payload, "name")
	if !ok { writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed.", []apiErrorDetail{{Field:"name", Code:"REQUIRED", Message:"Name is required."}}, requestID); return }
	message, ok := parseStringField(payload, "message")
	if !ok { writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed.", []apiErrorDetail{{Field:"message", Code:"REQUIRED", Message:"Message is required."}}, requestID); return }
	name = strings.TrimSpace(name); message = strings.TrimSpace(message)
	details := []apiErrorDetail{}
	if name == "" { details = append(details, apiErrorDetail{Field:"name", Code:"REQUIRED", Message:"Name is required."}) } else if utf8.RuneCountInString(name) > 80 { details = append(details, apiErrorDetail{Field:"name", Code:"TOO_LONG", Message:"Name must be 1 to 80 characters."}) }
	if message == "" { details = append(details, apiErrorDetail{Field:"message", Code:"REQUIRED", Message:"Message is required."}) } else if utf8.RuneCountInString(message) > 240 { details = append(details, apiErrorDetail{Field:"message", Code:"TOO_LONG", Message:"Message must be 1 to 240 characters."}) }
	if len(details) > 0 { writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed.", details, requestID); return }
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second); defer cancel()
	var row greeting
	if err := s.db.QueryRow(ctx, `INSERT INTO greetings (name, message) VALUES ($1, $2) RETURNING id, name, message, created_at`, name, message).Scan(&row.ID, &row.Name, &row.Message, &row.CreatedAt); err != nil { writeAPIError(w, http.StatusServiceUnavailable, "UNAVAILABLE", "Database unavailable.", nil, requestID); return }
	w.Header().Set("Location", "/v1/greetings/"+row.ID)
	writeJSON(w, http.StatusCreated, row)
}

func (s server) listGreetings(w http.ResponseWriter, r *http.Request) {
	requestID := apiRequestID(r)
	addRequestID(w, requestID)
	limit := 20
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed.", []apiErrorDetail{{Field: "limit", Code: "INVALID", Message: "Limit must be between 1 and 100."}}, requestID)
			return
		}
		limit = n
	}

	var cursor cursorPayload
	if raw := r.URL.Query().Get("cursor"); raw != "" {
		if err := decodeCursor(raw, &cursor); err != nil {
			writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed.", []apiErrorDetail{{Field: "cursor", Code: "INVALID", Message: "Cursor is malformed."}}, requestID)
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	query := `SELECT id::text, name, message, created_at FROM greetings`
	args := []any{limit + 1}
	if cursor.ID != "" {
		query += ` WHERE (created_at, id) < ($2::timestamptz, $3::bigint)`
		args = append(args, cursor.CreatedAt, cursor.ID)
	}
	query += ` ORDER BY created_at DESC, id DESC LIMIT $1`

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, "UNAVAILABLE", "Database unavailable.", nil, requestID)
		return
	}
	defer rows.Close()

	items := make([]greeting, 0, limit)
	for rows.Next() {
		var row greeting
		if err := rows.Scan(&row.ID, &row.Name, &row.Message, &row.CreatedAt); err != nil {
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL", "Unexpected error.", nil, requestID)
			return
		}
		items = append(items, row)
	}
	if err := rows.Err(); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "INTERNAL", "Unexpected error.", nil, requestID)
		return
	}

	hasMore := len(items) > limit
	var nextCursor *string
	if hasMore {
		items = items[:limit]
		c := encodeCursor(items[len(items)-1])
		nextCursor = &c
	}

	writeJSON(w, http.StatusOK, greetingsResponse{Greetings: items, NextCursor: nextCursor, HasMore: hasMore})
}

func parseStringField(payload map[string]json.RawMessage, key string) (string, bool) {
	raw, ok := payload[key]
	if !ok {
		return "", false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
}

func decodeCursor(raw string, out *cursorPayload) error {
	data, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func encodeCursor(row greeting) string {
	data, _ := json.Marshal(cursorPayload{CreatedAt: row.CreatedAt.UTC().Format(time.RFC3339), ID: row.ID})
	return base64.RawURLEncoding.EncodeToString(data)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}


func writeAPIError(w http.ResponseWriter, status int, code, message string, details []apiErrorDetail, requestID string) { var body apiError; body.Error.Code = code; body.Error.Message = message; body.Error.Details = details; body.Error.RequestID = requestID; writeJSON(w, status, body) }

func apiRequestID(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("X-Request-Id")); v != "" {
		return v
	}
	var b [12]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func addRequestID(w http.ResponseWriter, id string) {
	if id != "" {
		w.Header().Set("X-Request-Id", id)
	}
}


func applyMigrations(ctx context.Context, db *pgxpool.Pool) error {
	if _, err := db.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil { return err }
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil { return err }
	names := make([]string, 0, len(entries))
	for _, entry := range entries { if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") { names = append(names, entry.Name()) } }
	sort.Strings(names)
	for _, name := range names { if err := applyMigration(ctx, db, name); err != nil { return fmt.Errorf("%s: %w", name, err) } }
	return nil
}

func applyMigration(ctx context.Context, db *pgxpool.Pool, name string) error { tx, err := db.Begin(ctx); if err != nil { return err }; defer tx.Rollback(ctx); var exists bool; if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name = $1)`, name).Scan(&exists); err != nil { return err }; if exists { return tx.Commit(ctx) }; sqlBytes, err := fs.ReadFile(migrations.FS, name); if err != nil { return err }; if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil { return err }; if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1)`, name); err != nil { return err }; return tx.Commit(ctx) }

