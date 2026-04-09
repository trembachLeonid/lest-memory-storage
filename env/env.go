package env

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvVar string

func (env EnvVar) Get() string {
	return os.Getenv(string(env))
}

const (
	Env        EnvVar = "ENV"
	ShardCount EnvVar = "SHARD_COUNT"
)

func Load() error {
	return godotenv.Load(".env")
}
