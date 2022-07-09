package database

import (
	"context"

	"TwitchBot/config"
)

type DbServiceIFace interface {
	InitDB(ctx context.Context, config *config.DBConfig) error
	FindCommand(ctx context.Context, channel, command string) (string, bool, error)
	FindDuelCommand(ctx context.Context, channel string) (string, bool, error)
	FindDuelUser(ctx context.Context, username string) (*DuelStats, error)
	addNewUser(ctx context.Context, username string) error
}
