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
	Log            logger.LoggerConf                  `toml:"log"`
	CommonConf     CommonConfig                       `toml:"common"`
	ServerConf     ServerConfig                       `toml:"server"`
	TopicConf      map[string]map[string]MQConfConfig `toml:"topic"`
	Mysql          map[string]MySqlConfig             `toml:"mysql"`
	QuotaConf      QuotaConfig                        `toml:"quota"`
	RedisConf      RedisConfig                        `toml:"redis"`
	InterfaceHost  InterfaceHostConfig                `toml:"interface_host"`
	MultiEnvHost   map[string]MultiEnvHostConfig      `toml:"multi_env_host"`
	InterfaceInfo  map[string]InterfaceInfoConfig     `toml:"interface_info"`
	SkywalkingConf SkyWalkingConfig                   `toml:"skywalking"`
	RongCloudConf  RongCloudConfig                    `toml:"rongcloud"`
	KafkaConf      map[string]KafkaConfig             `toml:"kafka"`
	PulsarConf     map[string]PulsarConf              `toml:"pulsar"`
	PulsarTopic    map[string]PulsarTopic             `toml:"pulsar_topic"`
	Invite         Invite                             `toml:"invite"`
	ShortUrl       ShortUrl                           `toml:"short_url"`
}

type PulsarConf struct {
	Alias             string `toml:"alias"`
	Url               string `toml:"url"`
	OperationTimeout  int64  `toml:"operation_timeout"`
	ConnectionTimeout int64  `toml:"connection_timeout"`
	Token             string `toml:"token"`
}

type PulsarTopic struct {
	Topic         string `toml:"topic"`
	SubscribeName string `toml:"subscribe_name"`
}

type CommonConfig struct {
	CrashLogPath     string `toml:"crash_log_path"`
	RpcWorkPoolCount int64  `toml:"consumer_work_pool_count"`
	Env              string `toml:"env"`
	ServerName       string `toml:"server_name"`
}

type QuotaConfig struct {
	Enable        bool   `toml:"enable"`
	UseQConf      bool   `toml:"use_qconf"`
	QuotaQConfKey string `toml:"quota_qconf_key"`
	QuotaConfPath string `toml:"quota_conf_path"`
}

type ServerConfig struct {
	GPort    int `toml:"gport"`
	MPort    int `toml:"mport"`
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

type MQConfConfig struct {
	Topics        []string `toml:"topics"`
	Brokers       []string `toml:"brokers"`
	ConsumerGroup string   `toml:"group"`
}

type RedisConfig struct {
	MaxRetry     int    `toml:"max_retry"`
	Addr         string `toml:"addr"`
	PoolSize     int    `toml:"pool_size"`
	MinIdel      int    `toml:"min_idel"`
	IdelTimeout  int    `toml:"idle_timeout"`
	DialTimeout  int    `toml:"dial_timeout"`
	WriteTimeout int    `toml:"write_timeout"`
	ReadTimeout  int    `toml:"read_timeout"`
	Env          string `toml:"env"`
}

type InterfaceHostConfig map[string]string
type MultiEnvHostConfig map[string]string

type InterfaceInfoConfig struct {
	InterfacePath string `toml:"interface_path"`
	Timeout       int    `toml:"timeout"`
}

type SkyWalkingConfig struct {
	ServerName    string  `toml:"server_name"`
	LogPath       string  `toml:"log_path"`
	LogLinkPath   string  `toml:"log_link_path"`
	SamplerRate   float64 `toml:"sampler_rate"`
	Size          int     `toml:"size"`
	RotationCount uint    `toml:"rotation_count"`
}

type RongCloudConfig struct {
	AppKey    string `toml:"app_key"`
	AppSecret string `toml:"app_secret"`
	AppUri    string `toml:"app_uri"`
}

type KafkaConfig struct {
	BrokerList    string `toml:"broker_list"`
	Topic         string `toml:"topic"`
	GroupId       string `toml:"group_id"`
	ConsumerCount int    `toml:"consumer_count"`
}

type Invite struct {
	Url    string `toml:"url"`
	AesKey string `toml:"aes_key"`
}

type ShortUrl struct {
	Host string `toml:"host"`
	Path string `toml:"path"`
}
