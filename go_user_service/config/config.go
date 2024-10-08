package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"
)

type Config struct {
	ServiceName string
	Environment string // debug, test, release
	Version     string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword int
	PostgresDatabase string
	RedisHost        string
	RedisPort        string
	RedisPassword    string

	ContentServiceAddress string

	PostgresMaxConnections int32
}

// Load ...
func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.ServiceName = cast.ToString(getOrReturnDefaultValue("SERVICE_NAME", "go_user_service"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))
	config.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))

	config.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "user_db"))
	config.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
	config.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "mirodil"))
	config.PostgresPassword = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PASSWORD", 1212))
	config.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "to_do_list_user"))
	config.PostgresMaxConnections = cast.ToInt32(getOrReturnDefaultValue("POSTGRES_MAX_CONNECTIONS", 30))
	config.RedisHost = "redis" //cast.ToString(getOrReturnDefault("REDIS_HOST", "127.0.0.1"))
	config.RedisPort = "6379"  //cast.ToString(getOrReturnDefault("REDIS_PORT", "6379"))
	config.RedisPassword = cast.ToString(getOrReturnDefaultValue("REDIS_PASSWORD", "password"))

	config.ContentServiceAddress = cast.ToString(getOrReturnDefaultValue("USER_SERVICE_HOST", "go_user_service:8081"))
	// config.ContentGRPCPort = cast.ToString(getOrReturnDefaultValue("CONTENT_GRPC_PORT", "8081"))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
