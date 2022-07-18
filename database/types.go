package database

import (
	"context"

	"TwitchBot/config"
)

type DbServiceIFace interface {
	InitDB(ctx context.Context, config *config.DBConfig) error
	FindCommonCommand(ctx context.Context, channel, command string) (string, bool, error)
	FindDuelCommand(ctx context.Context, channel string) (string, bool, error)
	FindDuelUser(ctx context.Context, username string) (*DuelStats, error)
	addNewUser(ctx context.Context, username string) error
	GetDuelFinishCommand(ctx context.Context, channel string) (string, error)
	RefreshDuelStats(ctx context.Context, winner, loser string) error
	addNewDuelInStats(ctx context.Context, username string) error
	addNewDuelWin(ctx context.Context, username string) error
}
