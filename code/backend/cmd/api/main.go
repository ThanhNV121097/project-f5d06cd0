package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type server struct{ db *pgxpool.Pool }

type greeting struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type apiError struct{ Error errorBody `json:"error"` }

type errorBody struct {
	Code      string       `json:"code"`
	Message   string       `json:"message"`
	Details   []fieldError `json:"details,omitempty"`
	RequestID string       `json:"request_id"`
}

type fieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type helloResponse struct{ Message string `json:"message"` }

type createGreetingRequest struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type listResponse struct {
	Greetings  []greeting `json:"greetings"`
	NextCursor *string    `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}

type cursorPayload struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
}

const (
	maxRequestBody = 16 << 10
	maxNameLen     = 80
	maxMessageLen  = 240
)

type requestIDKey struct{}

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
	mux.HandleFunc("GET /v1/hello", srv.getHello)
	mux.HandleFunc("POST /v1/greetings", srv.createGreeting)
	mux.HandleFunc("GET /v1/greetings", srv.listGreetings)

	addr := ":" + listenPort()
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, requestIDMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func listenPort() string {
	if port := os.Getenv("PORT"); port != "" { return port }
	if port := os.Getenv("APP_PORT"); port != "" { return port }
	return "8080"
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" { id = fmt.Sprintf("req-%d", time.Now().UnixNano()) }
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id)))
	})
}

func requestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey{}).(string); ok { return v }
	return ""
}

func (s server) healthz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.db.Ping(ctx); err != nil { http.Error(w, "database unavailable", http.StatusServiceUnavailable); return }
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s server) apiHealth(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, map[string]string{"status": "ok"}) }

func (s server) getHello(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if utf8.RuneCountInString(name) > maxNameLen {
		writeAPIError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Name too long.", []fieldError{{Field: "name", Code: "TOO_LONG", Message: "Name must be 80 characters or fewer."}})
		return
	}
	if name == "" { name = "World" }
	writeJSON(w, http.StatusOK, helloResponse{Message: "Hello, " + name + "!"})
}

func (s server) createGreeting(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.HasPrefix(ct, "application/json") {
		writeAPIError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Content type must be JSON.", nil)
		return
	}
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req createGreetingRequest
	if err := dec.Decode(&req); err != nil {
		writeAPIError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request body must be valid JSON.", nil)
		return
	}
	if dec.More() {
		writeAPIError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request body must be a single JSON object.", nil)
		return
	}
	name, nameErr := validateText(req.Name, "name", maxNameLen)
	message, msgErr := validateText(req.Message, "message", maxMessageLen)
	if nameErr != nil || msgErr != nil {
		details := make([]fieldError, 0, 2)
		if nameErr != nil { details = append(details, *nameErr) }
		if msgErr != nil { details = append(details, *msgErr) }
		writeAPIError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Greeting fields are invalid.", details)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	var g greeting
	if err := s.db.QueryRow(ctx, `INSERT INTO greetings (name, message) VALUES ($1, $2) RETURNING id::text, name, message, created_at`, name, message).Scan(&g.ID, &g.Name, &g.Message, &g.CreatedAt); err != nil {
		respondDBError(w, r, err)
		return
	}
	w.Header().Set("Location", "/v1/greetings/"+g.ID)
	writeJSON(w, http.StatusCreated, g)
}

func (s server) listGreetings(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if raw := r.URL.Query().Get("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 || v > 100 {
			writeAPIError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Limit is invalid.", []fieldError{{Field: "limit", Code: "INVALID", Message: "Limit must be an integer from 1 to 100."}})
			return
		}
		limit = v
	}

	args := []any{}
	query := `SELECT id::text, name, message, created_at FROM greetings`
	if raw := r.URL.Query().Get("cursor"); raw != "" {
		payload, err := decodeCursor(raw)
		if err != nil { writeAPIError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Cursor is invalid.", []fieldError{{Field: "cursor", Code: "INVALID", Message: "Cursor must be a valid page token."}}); return }
		createdAt, err := time.Parse(time.RFC3339Nano, payload.CreatedAt)
		if err != nil { writeAPIError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Cursor is invalid.", []fieldError{{Field: "cursor", Code: "INVALID", Message: "Cursor must be a valid page token."}}); return }
		id, err := strconv.ParseInt(payload.ID, 10, 64)
		if err != nil { writeAPIError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Cursor is invalid.", []fieldError{{Field: "cursor", Code: "INVALID", Message: "Cursor must be a valid page token."}}); return }
		query += ` WHERE (created_at, id) < ($1, $2)`
		args = append(args, createdAt, id)
	}
	query += ` ORDER BY created_at DESC, id DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit+1)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil { respondDBError(w, r, err); return }
	defer rows.Close()

	items := make([]greeting, 0, limit)
	for rows.Next() {
		var row greeting
		if err := rows.Scan(&row.ID, &row.Name, &row.Message, &row.CreatedAt); err != nil { respondDBError(w, r, err); return }
		items = append(items, row)
	}
	if err := rows.Err(); err != nil { respondDBError(w, r, err); return }

	hasMore := len(items) > limit
	if hasMore { items = items[:limit] }
	var nextCursor *string
	if hasMore && len(items) > 0 {
		last := items[len(items)-1]
		encoded := encodeCursor(cursorPayload{CreatedAt: last.CreatedAt, ID: last.ID})
		nextCursor = &encoded
	}
	writeJSON(w, http.StatusOK, listResponse{Greetings: items, NextCursor: nextCursor, HasMore: hasMore})
}

func validateText(raw, field string, maxLen int) (string, *fieldError) {
	value := strings.TrimSpace(raw)
	label := strings.ToUpper(field[:1]) + field[1:]
	if value == "" { return "", &fieldError{Field: field, Code: "REQUIRED", Message: label + " is required."} }
	if utf8.RuneCountInString(value) > maxLen { return "", &fieldError{Field: field, Code: "TOO_LONG", Message: fmt.Sprintf("%s must be %d characters or fewer.", label, maxLen)} }
	return value, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeAPIError(w http.ResponseWriter, r *http.Request, status int, code, message string, details []fieldError) {
	writeJSON(w, status, apiError{Error: errorBody{Code: code, Message: message, Details: details, RequestID: requestIDFromContext(r.Context())}})
}

func respondDBError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) { writeAPIError(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "Service temporarily unavailable.", nil); return }
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) { writeAPIError(w, r, http.StatusInternalServerError, "INTERNAL", "Unexpected database failure.", nil); return }
	if errors.Is(err, pgx.ErrNoRows) { writeAPIError(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "Service temporarily unavailable.", nil); return }
	writeAPIError(w, r, http.StatusInternalServerError, "INTERNAL", "Unexpected server failure.", nil)
}

func decodeCursor(raw string) (cursorPayload, error) {
	var payload cursorPayload
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil { return payload, err }
	if err := json.Unmarshal(b, &payload); err != nil { return payload, err }
	return payload, nil
}

func encodeCursor(payload cursorPayload) string {
	b, _ := json.Marshal(payload)
	return base64.RawURLEncoding.EncodeToString(b)
}

func applyMigrations(ctx context.Context, db *pgxpool.Pool) error {
	if _, err := db.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil { return err }
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil { return err }
	names := make([]string, 0, len(entries))
	for _, entry := range entries { if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") { names = append(names, entry.Name()) } }
	sort.Strings(names)
	for _, name := range names {
		if err := applyMigration(ctx, db, name); err != nil { return fmt.Errorf("%s: %w", name, err) }
	}
	return nil
}

func applyMigration(ctx context.Context, db *pgxpool.Pool, name string) error {
	tx, err := db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name = $1)`, name).Scan(&exists); err != nil { return err }
	if exists { return tx.Commit(ctx) }
	sqlBytes, err := fs.ReadFile(migrations.FS, name)
	if err != nil { return err }
	if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil { return err }
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1)`, name); err != nil { return err }
	return tx.Commit(ctx)
}
