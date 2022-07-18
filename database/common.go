package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type DBClient struct {
	db *pgx.Conn
}

func (db *DBClient) FindCommonCommand(ctx context.Context, channel, command string) (string, bool, error) {
	var answer string

	err := db.db.QueryRow(ctx,
		`SELECT answer 
			FROM commands
			WHERE channel = $1 AND command = $2`, channel, command).Scan(&answer)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		} else {
			return "", false, err
		}
	}

	return answer, true, nil
}
