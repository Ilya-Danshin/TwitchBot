package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Pass     string
}

type Config struct {
	DBConf *DBConfig
	Users  []*User
}

type User struct {
	Name    string   `json:"name"`
	Modules []string `json:"modules"`
}

func InitConfig() error {
	pathToEnvFile, ok := os.LookupEnv("PATH_TO_ENV_FILE")
	if !ok {
		return fmt.Errorf("env PATH_TO_ENV_FILE is not set")
	}

	err := godotenv.Load(pathToEnvFile)
	if err != nil {
		return fmt.Errorf("can`t load env variables: %w", err)
	}

	return nil
}

func ParseConfig() (*Config, error) {
	dbConf, err := getDBConfig()
	if err != nil {
		return nil, err
	}

	users, err := getUsers()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBConf: dbConf,
		Users:  users,
	}, nil
}

func getDBConfig() (*DBConfig, error) {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return nil, fmt.Errorf("env DB_HOST is not set")
	}

	temp, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return nil, fmt.Errorf("env DB_PORT is not set")
	}

	port, err := strconv.Atoi(temp)
	if err != nil {
		return nil, err
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, fmt.Errorf("env DB_NAME is not set")
	}

	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		return nil, fmt.Errorf("env DB_USER is not set")
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("env DB_PASSWORD is not set")
	}

	return &DBConfig{
		Host:     host,
		Port:     port,
		Database: dbName,
		User:     user,
		Pass:     password,
	}, nil
}

func getUsers() ([]*User, error) {
	path, ok := os.LookupEnv("USERS_SETTINGS_PATH")
	if !ok {
		return nil, fmt.Errorf("env USERS_PATH is not set")
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var users []*User
	err = json.Unmarshal(content, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
