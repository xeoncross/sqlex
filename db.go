package sqlex

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL error code for duplicate rows
const ER_DUP_ENTRY = 1062

// ensure instance implements db
var _ db = (*DB)(nil)

type db interface {
	// database/sql "bare-metal"
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row

	// sqlx
	SelectContext(ctx context.Context, instances interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, instance interface{}, query string, args ...interface{}) error
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)

	// Helpers
	InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error)
	InsertContextPsql(ctx context.Context, query string, args ...interface{}) (int64, error)
	InsertIgnoreContext(ctx context.Context, query string, args ...interface{}) (bool, error)
	UpdateContext(ctx context.Context, query string, args ...interface{}) (int64, error)
	UpdateOneContext(ctx context.Context, query string, args ...interface{}) error
	GetOrCreate(ctx context.Context, dest interface{}, selectQuery, insertQuery string, args ...interface{}) error
}

type DB struct {
	*sqlx.DB
}

// Open is the same as sql.Open, but returns an *sqlex.DB instead.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (s *DB) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("insert: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insert: %w", err)
	}
	return id, nil
}

func (s *DB) InsertIgnoreContext(ctx context.Context, query string, args ...interface{}) (bool, error) {
	result, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("insert: %w", err)
	}
	num, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("insert: %w", err)
	}

	// Record already exists or new record(s)?
	return num != 0, nil
}

func (s *DB) UpdateContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("update: %w", err)
	}
	num, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update: %w", err)
	}
	return num, nil
}

func (s *DB) UpdateOneContext(ctx context.Context, query string, args ...interface{}) error {
	res, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if num != 1 {
		return fmt.Errorf("update: %d rows affected, expected 1", num)
	}
	return nil
}

// GetOrCreate record using the first argument as the unique lookup key for the record and the rest for the insert
func (s *DB) GetOrCreate(ctx context.Context, dest interface{}, selectQuery, insertQuery string, args ...interface{}) error {

	if len(args) == 0 {
		return errors.New("no arguments provided")
	}

	first := args[0]
	err := s.GetContext(ctx, dest, selectQuery, first)
	if err == nil {
		return nil
	}

	// Some actual database error?
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Try to create it, then load it
	_, err = s.InsertContext(ctx, insertQuery, args...)
	if err == nil {
		return s.GetContext(ctx, dest, selectQuery, first)
	}

	// MySQL specific error check
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		// Another process created this entry since we ran SELECT
		if mysqlError.Number == ER_DUP_ENTRY {
			return s.GetContext(ctx, dest, selectQuery, first)
		}
	}

	// TODO: Postgres specific error check

	return err
}

// InsertContextPsql inserts a single record reading RETURNING id
func (s *DB) InsertContextPsql(ctx context.Context, query string, args ...interface{}) (int64, error) {
	row := s.DB.QueryRowContext(ctx, query, args...)
	// TODO: GetContext for multi-value returns
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
