package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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

const (
	maxRequestBodyBytes = 16 << 10
	nameLimit           = 80
	messageLimit        = 240
	defaultListLimit    = 20
	maxListLimit        = 100
)

type server struct {
	db *pgxpool.Pool
}

type apiError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []fieldDetail `json:"details,omitempty"`
}

type fieldDetail struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type greeting struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type greetingListResponse struct {
	Greetings  []greeting `json:"greetings"`
	NextCursor *string    `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}

type cursorPayload struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
}

type createGreetingInput struct {
	Name    string
	Message string
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
	mux.HandleFunc("GET /api/health", srv.health)
	mux.HandleFunc("GET /api/hello", srv.hello)
	mux.HandleFunc("POST /api/greetings", srv.createGreeting)
	mux.HandleFunc("GET /api/greetings", srv.listGreetings)

	addr := ":" + listenPort()
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, requestLogger(mux)); err != nil {
		log.Fatal(err)
	}
}

func listenPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		return port
	}
	return "8080"
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		_ = start
	})
}

func (s server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s server) healthz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := s.db.Ping(ctx); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s server) hello(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if utf8.RuneCountInString(name) > nameLimit {
		writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Name is too long.", []fieldDetail{{Field: "name", Code: "TOO_LONG", Message: "Name must be 80 characters or fewer."}})
		return
	}
	if name == "" {
		name = "World"
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Hello, %s!", name)})
}

func (s server) createGreeting(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" || !strings.HasPrefix(contentType, "application/json") {
		writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Request body must be JSON.", nil)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxRequestBodyBytes+1))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Request body could not be read.", nil)
		return
	}
	if len(body) > maxRequestBodyBytes {
		writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Request body is too large.", nil)
		return
	}
	var input struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		if errors.Is(err, io.EOF) {
			writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Request body must be valid JSON.", nil)
			return
		}
		if strings.Contains(err.Error(), "cannot unmarshal") || strings.Contains(err.Error(), "unknown field") || strings.Contains(err.Error(), "invalid character") || strings.Contains(err.Error(), "number") {
			writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Request body must contain strings.", nil)
			return
		}
		writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Request body must be valid JSON.", nil)
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		writeAPIError(w, http.StatusBadRequest, "BAD_REQUEST", "Request body must be a single JSON object.", nil)
		return
	}

	name, message, details := validateGreetingInput(createGreetingInput{Name: input.Name, Message: input.Message})
	if len(details) > 0 {
		writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Greeting is invalid.", details)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var created greeting
	var createdAt time.Time
	err = s.db.QueryRow(ctx, `INSERT INTO greetings (name, message) VALUES ($1, $2) RETURNING id, name, message, created_at`, name, message).Scan(&created.ID, &created.Name, &created.Message, &createdAt)
	if err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, "UNAVAILABLE", "Database unavailable.", nil)
		return
	}
	created.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	w.Header().Set("Location", "/api/greetings/"+created.ID)
	writeJSON(w, http.StatusCreated, created)
}

func (s server) listGreetings(w http.ResponseWriter, r *http.Request) {
	limit, cursor, errResp := parseListParams(r)
	if errResp != nil {
		writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Query parameters are invalid.", errResp)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	query := `SELECT id, name, message, created_at FROM greetings`
	args := []any{}
	if cursor != nil {
		query += ` WHERE (created_at, id) < ($1, $2)`
		args = append(args, cursor.CreatedAt, cursor.ID)
	}
	query += ` ORDER BY created_at DESC, id DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit+1)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, "UNAVAILABLE", "Database unavailable.", nil)
		return
	}
	defer rows.Close()

	items := make([]greeting, 0, limit+1)
	for rows.Next() {
		var item greeting
		var createdAt time.Time
		if err := rows.Scan(&item.ID, &item.Name, &item.Message, &createdAt); err != nil {
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL", "Unexpected failure.", nil)
			return
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "INTERNAL", "Unexpected failure.", nil)
		return
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	var nextCursor *string
	if hasMore && len(items) > 0 {
		payload := cursorPayload{CreatedAt: items[len(items)-1].CreatedAt, ID: items[len(items)-1].ID}
		encoded, _ := encodeCursor(payload)
		nextCursor = &encoded
	}
	writeJSON(w, http.StatusOK, greetingListResponse{Greetings: items, NextCursor: nextCursor, HasMore: hasMore})
}

func validateGreetingInput(input createGreetingInput) (string, string, []fieldDetail) {
	name := strings.TrimSpace(input.Name)
	message := strings.TrimSpace(input.Message)
	var details []fieldDetail
	if name == "" {
		details = append(details, fieldDetail{Field: "name", Code: "REQUIRED", Message: "Name is required."})
	} else if utf8.RuneCountInString(name) > nameLimit {
		details = append(details, fieldDetail{Field: "name", Code: "TOO_LONG", Message: "Name must be 80 characters or fewer."})
	}
	if message == "" {
		details = append(details, fieldDetail{Field: "message", Code: "REQUIRED", Message: "Message is required."})
	} else if utf8.RuneCountInString(message) > messageLimit {
		details = append(details, fieldDetail{Field: "message", Code: "TOO_LONG", Message: "Message must be 240 characters or fewer."})
	}
	return name, message, details
}

func parseListParams(r *http.Request) (int, *cursorPayload, []fieldDetail) {
	limit := defaultListLimit
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > maxListLimit {
			return 0, nil, []fieldDetail{{Field: "limit", Code: "INVALID", Message: "Limit must be between 1 and 100."}}
		}
		limit = parsed
	}
	if raw := r.URL.Query().Get("cursor"); raw != "" {
		cursor, err := decodeCursor(raw)
		if err != nil {
			return 0, nil, []fieldDetail{{Field: "cursor", Code: "INVALID", Message: "Cursor is invalid."}}
		}
		return limit, cursor, nil
	}
	return limit, nil, nil
}

func encodeCursor(p cursorPayload) (string, error) {
	b, err := json.Marshal(p)
	if err != nil { return "", err }
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func decodeCursor(raw string) (*cursorPayload, error) {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil { return nil, err }
	var p cursorPayload
	if err := json.Unmarshal(b, &p); err != nil { return nil, err }
	if p.CreatedAt == "" || p.ID == "" { return nil, errors.New("cursor incomplete") }
	return &p, nil
}

func writeAPIError(w http.ResponseWriter, status int, code, message string, details []fieldDetail) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: apiError{Code: code, Message: message, Details: details}})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func applyMigrations(ctx context.Context, db *pgxpool.Pool) error {
	if _, err := db.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		return err
	}

	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		if err := applyMigration(ctx, db, name); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func applyMigration(ctx context.Context, db *pgxpool.Pool, name string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name = $1)`, name).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return tx.Commit(ctx)
	}

	sqlBytes, err := fs.ReadFile(migrations.FS, name)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1)`, name); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

