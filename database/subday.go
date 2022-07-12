package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func (db *DBClient) FindSubdayCommand(ctx context.Context, channel string) (string, bool, error) {
	var answer string

	err := db.db.QueryRow(ctx,
		`SELECT answer
			FROM subdays_command
			WHERE channel = $1`, channel).Scan(&answer)

	if err != nil {
		if err != pgx.ErrNoRows {
			return "", false, err
		}
		return "", false, nil
	}

	return answer, true, err
}

func (db *DBClient) AddNewSubdayOrder(ctx context.Context, channel, nickname, order string) error {
	err := db.db.QueryRow(ctx,
		`INSERT INTO subdays
			VALUES ($1, $2, $3) 
			ON CONFLICT (channel, nickname) 
			DO UPDATE 
			SET "order"=$3`, channel, nickname, order).Scan()

	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	return nil
}
