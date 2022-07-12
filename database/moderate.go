package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func (db *DBClient) FindModerateCommand(ctx context.Context, channel, command string) (string, bool, error) {
	var answer string
	err := db.db.QueryRow(ctx,
		`SELECT answer
			FROM moderate_commands
			WHERE channel = $1 AND command = $2`, channel, command).Scan(&answer)

	if err != nil {
		if err != pgx.ErrNoRows {
			return "", false, err
		}
		return "", false, nil
	}

	return answer, true, nil
}
