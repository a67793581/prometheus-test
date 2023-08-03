package environment

import (
	"os"
	"time"
)

const (
	EnvUnknown = "unknown"
)

func InitEnv() {
	EnvCfg = &EnvConfig{}
	EnvCfg.StartTime = time.Now()
}

func getOsEnv(key string) string {
	env := os.Getenv(key)
	if env != "" {
		return env
	}
	return EnvUnknown
}

type EnvConfig struct {
	StartTime time.Time
}

var EnvCfg *EnvConfig
