package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type NewEmail struct {
	Email string `json:"email"`
	Link  string `json:"link"`
}

func InitDB() error {
	// Create the schema and table if they don't exist
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	// Check if schema exists
	tsql := fmt.Sprintf("SELECT COUNT(*) FROM sys.schemas WHERE name = '%s'", schema)
	row := db.QueryRowContext(ctx, tsql)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		// Create schema
		tsql = fmt.Sprintf("CREATE SCHEMA %s;", schema)
		_, err = db.ExecContext(ctx, tsql)
		if err != nil {
			return err
		}
	}

	// Check if table exists
	tsql = fmt.Sprintf("SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s';", schema, tableName)
	row = db.QueryRowContext(ctx, tsql)
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		// Create table
		tsql = fmt.Sprintf("CREATE TABLE %s.%s (Id INT IDENTITY(1,1) PRIMARY KEY, Email NVARCHAR(250) NOT NULL, Link NVARCHAR(250) NOT NULL, CreatedAt DATETIME2 NOT NULL, UsedAt DATETIME2);", schema, tableName)
		_, err = db.ExecContext(ctx, tsql)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddEmailInfo(emailInfo NewEmail) (int64, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}
	log.Println(fmt.Sprintf("%s.%s", schema, tableName))

	tsql := fmt.Sprintf("INSERT INTO %s.%s (Email, Link, CreatedAt) VALUES (@Email, @Link, @CreatedAt);", schema, tableName)

	// Execute non-query with named parameters
	result, err := db.ExecContext(ctx, tsql, sql.Named("Email", emailInfo.Email), sql.Named("Link", emailInfo.Link), sql.Named("CreatedAt", time.Now()))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func UsedLink(link string) (int64, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("UPDATE %s.%s SET UsedAt = @UsedAt WHERE Link = @Link;", schema, tableName)

	// Execute non-query with named parameters
	result, err := db.ExecContext(ctx, tsql, sql.Named("UsedAt", time.Now()), sql.Named("Link", link))

	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
