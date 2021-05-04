# sqlex

Simple additional methods to make using https://github.com/jmoiron/sqlx easier for MySQL and Postgres databases. 

## Example

This library was written to reduce the amount of code that had to be typed when writing out 50+ database model methods when using auto-incrementing primary keys with PostgreSQL or MySQL.

Example sqlx (mysql) insert code:

```
func (s *MySQL) InsertUser(ctx context.Context, user *myapp.User) (int64, error) {
	query := `INSERT INTO user(
		email,
		password,
        created_at
	) values (?,?,UNIX_TIMESTAMP(now()))`
    res, err := s.DB.ExecContext(ctx, sql, user.Email, user.Password)

    if err != nil {
        return 0, fmt.Errorf("insertuser: %w", err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("insertuser: %w", err)
    }
	
	return id, nil
}
```

with sqlex we can write this much simpler

```
func (m *Model) InsertUser(ctx context.Context, user *myapp.User) (int64, error) {
	query := `INSERT INTO user(
		email,
		password,
        created_at
	) values (?,?,UNIX_TIMESTAMP(now()))`

	return m.DB.InsertContext(ctx, query, user.Email, user.Password)
}
```

## Usage

Simply open a sqlex connection the same as you would an sql or sqlx connection.

```go
db, err := sqlex.Open("mysql", "user:pass@tcp(localhost:3306)/myapp?collation=utf8mb4_unicode_ci&parseTime=true")
```

In addition to all sqlx methods, you now have access to the following helpers:

```
InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error)
InsertContextPsql(ctx context.Context, query string, args ...interface{}) (int64, error)
InsertIgnoreContext(ctx context.Context, query string, args ...interface{}) (bool, error)
UpdateContext(ctx context.Context, query string, args ...interface{}) (int64, error)
UpdateOneContext(ctx context.Context, query string, args ...interface{}) error
GetOrCreate(ctx context.Context, dest interface{}, selectQuery, insertQuery string, args ...interface{}) error
```

Furthermore, two additional query wrappers are also available: `ErrorWrapper` and `Logger`.

```go
db, err := sqlex.Open(...)

Wrap errors with store caller name
db = &sqlex.ErrorWrapper{db}
```

## Notes

Don't forget to load the actual driver in your application

```
package main

import (
    ...
    _ "github.com/go-sql-driver/mysql"
```
