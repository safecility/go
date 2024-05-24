package setup

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"os"
)

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Key  string `json:"key"`
}

func (r *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
func (r *RedisConfig) NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     r.Address(),
		Password: r.Key,
		DB:       0, // use default DB
	})
}

type MySQLConfig struct {
	//dont json in password
	Password               string
	Username               string `json:"user"`
	Host                   string `json:"host"`
	Port                   int    `json:"port"`
	Database               string `json:"database"`
	InstanceConnectionName string `json:"instanceConnectionName"`
}

// connectionString returns a connection string suitable for sql.Open.
func (c MySQLConfig) connectionString(databaseName string) string {
	var cred string
	// [username[:password]@]
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	//"%s:%s@unix(/%s/%s)/%s?parseTime=true"  dbUser, dbPwd, socketDir, instanceConnectionName, dbName
	if c.InstanceConnectionName != "" {
		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}
		log.Debug().Str("instance", c.InstanceConnectionName).Msg("using instance connection")
		return fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true",
			c.Username, c.Password, socketDir, c.InstanceConnectionName, databaseName)
	}
	return fmt.Sprintf("%stcp([%s]:%d)/%s?parseTime=true", cred, c.Host, c.Port, databaseName)
}

// NewSafecilitySql creates a new MySQL server.
func NewSafecilitySql(config MySQLConfig) (*sql.DB, error) {
	//we need to load the mysql driver by reference
	_ = mysql.ErrBusyBuffer
	c := config.connectionString(config.Database)
	db, err := sql.Open("mysql", c)
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
	}

	return db, nil
}
