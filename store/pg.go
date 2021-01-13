package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	ini "github.com/vaughan0/go-ini"
	"github.com/zoommix/fasthttp_template/utils"
)

// DBConfig ..
type DBConfig struct {
	Name     string
	Host     string
	User     string
	Password string
}

const (
	maxConns = 4
	minConns = 2

	truncateSQL     = "TRUNCATE %s CASCADE"
	testENV         = "test"
	configPath      = "config/%s.conf"
	testDatabaseURL = "host=localhost dbname=fasthttp_template_test sslmode=disable"
)

// DB ...
var DB *pgxpool.Pool
var dbConfig DBConfig

func init() {
	var err error

	dbURL := fetchDBURL()

	config, err := pgxpool.ParseConfig(dbURL)

	config.MaxConns = maxConns
	config.MinConns = minConns
	config.ConnConfig.Logger = &utils.PgxLogger{}

	if err != nil {
		panic(err)
	}

	DB, err = pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}
}

// Close ...
func Close() {
	DB.Close()
}

func fetchDBURL() string {
	var url string

	if utils.GetENV() == testENV {
		url = testDatabaseURL
	} else {
		dbConf := fetchDBConfig()

		url = fmt.Sprintf(
			"dbname=%s host=%s user=%s password=%s sslmode=disable",
			dbConf.Name,
			dbConf.Host,
			dbConf.User,
			dbConf.Password,
		)
	}

	return url
}

func fetchDBConfig() DBConfig {
	env := utils.GetENV()
	path := fmt.Sprintf(configPath, env)
	file, err := ini.LoadFile(path)

	if err != nil {
		panic(err)
	}

	dbConfig.Host, _ = file.Get("database", "host")
	dbConfig.Name, _ = file.Get("database", "database")
	dbConfig.Password, _ = file.Get("database", "password")
	dbConfig.User, _ = file.Get("database", "user")

	return dbConfig
}

// TearDown truncates specified database table
func TearDown(tables ...string) {
	if utils.GetENV() == testENV {
		if len(tables) > 0 {
			_, err := DB.Exec(context.Background(), fmt.Sprintf(truncateSQL, strings.Join(tables, ",")))

			if err != nil {
				panic(err)
			}
		}
	}
}
