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
			float64(s.wins)/float64(s.duels)*100)
	}
	return "(duels: 0, wins: 0, winrate: 0.0%)"
}

func (db *DBClient) FindDuelCommand(ctx context.Context, channel string) (string, bool, error) {
	var answer string

	err := db.db.QueryRow(ctx,
		`SELECT introduction 
			FROM duel_commands
			WHERE channel = $1`, channel).Scan(&answer)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	return answer, true, nil
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
			VALUES ($1)`, username)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}

func (db *DBClient) GetDuelFinishCommand(ctx context.Context, channel string) (string, error) {
	var answer string
	err := db.db.QueryRow(ctx,
		`SELECT conclusion 
			FROM duel_commands
			WHERE channel=$1`, channel).Scan(&answer)

	if err != nil {
		return "", err
	}

	return answer, nil
}

func (db *DBClient) RefreshDuelStats(ctx context.Context, winner, loser string) error {
	err := db.addNewDuelInStats(ctx, winner)
	if err != nil {
		return err
	}

	err = db.addNewDuelInStats(ctx, loser)
	if err != nil {
		return err
	}

	err = db.addNewDuelWin(ctx, winner)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBClient) addNewDuelInStats(ctx context.Context, username string) error {
	err := db.db.QueryRow(ctx,
		`UPDATE duel_statistics
			SET duels = duels + 1
			WHERE nickname = $1`, username).Scan()
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	return nil
}

func (db *DBClient) addNewDuelWin(ctx context.Context, username string) error {
	err := db.db.QueryRow(ctx,
		`UPDATE duel_statistics
			SET wins = wins + 1
			WHERE nickname = $1`, username).Scan()
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	return nil
}
