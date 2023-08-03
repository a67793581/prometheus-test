package config

import (
	"log"

	"github.com/BurntSushi/toml"

	"prometheus-test/lib/logger"
)

var Cfg Config

func InitConfig(filePath string) error {
	if _, err := toml.DecodeFile(filePath, &Cfg); err != nil {
		return err
	}
	log.Printf("DisPatcher_Config=%+v", Cfg)

	return nil
}

type Config struct {
	Log        logger.LoggerConf      `toml:"log"`
	CommonConf CommonConfig           `toml:"common"`
	ServerConf ServerConfig           `toml:"server"`
	Mysql      map[string]MySqlConfig `toml:"mysql"`
}

type CommonConfig struct {
	CrashLogPath string `toml:"crash_log_path"`
	Env          string `toml:"env"`
	ServerName   string `toml:"server_name"`
}

type ServerConfig struct {
	GPort    int `toml:"gport"`
	WTimeout int `toml:"wTimeout"`
	RTimeout int `toml:"wTimeout"`
}

type MySqlConfig struct {
	DBName          string `toml:"db_name"`
	Host            string `toml:"host"`
	ReadHost        string `toml:"read_host"`
	Port            int    `toml:"port"`
	User            string `toml:"user"`
	Passwd          string `toml:"passwd"`
	ConnTimeout     string `toml:"conn_timeout"`
	ReadTimeout     string `toml:"read_timeout"`
	WriteTimeout    string `toml:"write_timeout"`
	MaxConnNum      int    `toml:"max_conn_num"`
	MaxIdleConnNum  int    `toml:"max_idle_conn_num"`
	MaxConnLifeTime int    `toml:"max_conn_life_time"`
	LogLevel        int    `toml:"log_level"`
}
