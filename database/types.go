package database

import (
	"TwitchBot/config"
	"context"
)

type DbServiceIFace interface {
	InitDB(ctx context.Context, config *config.DBConfig) error
}
