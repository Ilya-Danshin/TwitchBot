package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type DuelStats struct {
	duels int
	wins  int
}

func (s *DuelStats) String() string {
	if s.duels != 0 {
		return fmt.Sprintf("(duels: %d, wins: %d, winrate: %.1f%%)", s.duels, s.wins,
			float64(s.wins)/float64(s.duels))
	}
	return "(duels: 0, wins: 0, winrate: 0.0%)"
}

func (db *DBClient) FindDuelCommand(ctx context.Context, channel string) (string, bool, error) {
	return db.FindCommand(ctx, channel, "дуэль")
}

func (db *DBClient) FindDuelUser(ctx context.Context, username string) (*DuelStats, error) {
	var stats DuelStats
	err := db.db.QueryRow(ctx,
		`SELECT duels, wins
			FROM duel_statistics
			WHERE nickname = $1`, username).Scan(&stats.duels, &stats.wins)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = db.addNewUser(ctx, username)
			if err != nil {
				return nil, err
			}
			return &DuelStats{duels: 0,
				wins: 0}, nil
		}
		return nil, err
	}

	return &stats, nil
}

func (db *DBClient) addNewUser(ctx context.Context, username string) error {
	rows, err := db.db.Query(ctx,
		`INSERT INTO duel_statistics
			VALUES ($1, 0, 0)`, username)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
