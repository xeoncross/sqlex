package sqlex

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
)

//
// Error Wrapper
//

var _ db = (*ErrorWrapper)(nil)

// ErrorWrapper for database queries
type ErrorWrapper struct {
	*DB
}

func (l *ErrorWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	r, err := l.DB.ExecContext(ctx, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return r, err
}

func (l *ErrorWrapper) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	r, err := l.DB.PrepareContext(ctx, query)
	if err != nil {
		err = wrapError(err)
	}
	return r, err
}

func (l *ErrorWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	r, err := l.DB.QueryContext(ctx, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return r, err
}

func (l *ErrorWrapper) SelectContext(ctx context.Context, DBs interface{}, query string, args ...interface{}) error {
	err := l.DB.SelectContext(ctx, DBs, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return err
}

func (l *ErrorWrapper) GetContext(ctx context.Context, DB interface{}, query string, args ...interface{}) error {
	err := l.DB.GetContext(ctx, DB, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return err
}

func (l *ErrorWrapper) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	r, err := l.DB.InsertContext(ctx, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return r, err
}

func (l *ErrorWrapper) InsertIgnoreContext(ctx context.Context, query string, args ...interface{}) (bool, error) {
	r, err := l.DB.InsertIgnoreContext(ctx, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return r, err
}

func (l *ErrorWrapper) UpdateContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	r, err := l.DB.UpdateContext(ctx, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return r, err
}

func (l *ErrorWrapper) UpdateOneContext(ctx context.Context, query string, args ...interface{}) error {
	err := l.DB.UpdateOneContext(ctx, query, args...)
	if err != nil {
		err = wrapError(err)
	}
	return err
}

// We do our best to avoid the overhead of calling this in the above methods
func wrapError(err error) error {

	fpcs := make([]uintptr, 1)
	// Skip 3 levels to get the caller
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return err
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return err
	}

	// file name and line number..?
	// caller.FileLine(fpcs[0]-1)

	return fmt.Errorf("%s: %w", filepath.Base(caller.Name()), err)
}
