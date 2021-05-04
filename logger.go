package sqlex

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"
)

//
// Query Logger
//

var _ db = (*Logger)(nil)

// Logger for database queries
type Logger struct {
	*DB
	Log *log.Logger
}

// ExecContext shell
func (l *Logger) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start), cleanWhitespace(query), args)
	}()
	return l.DB.ExecContext(ctx, query, args...)
}

// PrepareContext shell
func (l *Logger) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query))
	}()
	return l.DB.PrepareContext(ctx, query)
}

// QueryContext shell
func (l *Logger) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext shell
func (l *Logger) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.QueryRowContext(ctx, query, args...)
}

func (l *Logger) SelectContext(ctx context.Context, DBs interface{}, query string, args ...interface{}) error {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.SelectContext(ctx, DBs, query, args...)
}

func (l *Logger) GetContext(ctx context.Context, DB interface{}, query string, args ...interface{}) error {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.GetContext(ctx, DB, query, args...)
}

func (l *Logger) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.InsertContext(ctx, query, args...)
}

func (l *Logger) InsertIgnoreContext(ctx context.Context, query string, args ...interface{}) (bool, error) {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.InsertIgnoreContext(ctx, query, args...)
}

func (l *Logger) UpdateContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.UpdateContext(ctx, query, args...)
}

func (l *Logger) UpdateOneContext(ctx context.Context, query string, args ...interface{}) error {
	start := time.Now()
	defer func() {
		l.Log.Printf("%s: %s: %v\n", time.Since(start).Round(time.Millisecond), cleanWhitespace(query), args)
	}()
	return l.DB.UpdateOneContext(ctx, query, args...)
}

func cleanWhitespace(str string) string {
	return strings.Join(strings.Fields(str), " ")
}
