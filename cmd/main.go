package main

import (
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

	fmt.Printf("host: %s\nport: %d\ndatabese name: %s\nuser: %s\npass:%s\n", cfg.DBConf.Host, cfg.DBConf.Port,
		cfg.DBConf.Database, cfg.DBConf.User, cfg.DBConf.Pass)

	for _, user := range cfg.Users {
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Modules: %s\n", user.Modules)
	}
}
