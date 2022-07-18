package database

import (
	"context"
	"fmt"

	"TwitchBot/config"

	"github.com/jackc/pgx/v4"
)

var DB DBClient

func (db *DBClient) InitDB(ctx context.Context, config *config.DBConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.Host, config.User,
		config.Pass, config.Database, config.Port)

	pgxConnection, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	DB.db = pgxConnection
	return nil
}
