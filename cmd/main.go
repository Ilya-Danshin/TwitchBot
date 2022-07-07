package main

import (
	"TwitchBot/database"
	"TwitchBot/internal/bot"
	"context"
	"fmt"

	"TwitchBot/config"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		fmt.Println("initial config error: ", err.Error())
		return
	}

	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Println("parse config error: ", err.Error())
	}

	err = database.DB.InitDB(context.Background(), cfg.DBConf)
	if err != nil {
		fmt.Println("database initialization error: ", err.Error())
	}

	bot.InitBot(cfg.Users, cfg.Bot)
	bot.LoopBot()
}
