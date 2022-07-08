package database

import (
	"context"
	"fmt"

	"TwitchBot/config"

	"github.com/jackc/pgx/v4"
)

type DBClient struct {
	db *pgx.Conn
}

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

func (db *DBClient) FindCommand(ctx context.Context, channel, command string) (string, error) {
	var answer string

	err := db.db.QueryRow(ctx,
		`SELECT answer 
			FROM commands
			WHERE channel = $1 AND command = $2`, channel, command).Scan(&answer)
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", err
		}
	}

	return answer, nil
}
