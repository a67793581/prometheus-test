package drivers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"prometheus-test/config"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	ErrNotFoundEngine = errors.New("not found sql engine")
	engineManager     = make(map[string]*gorm.DB)
)

func InitMysql() error {
	for dbName, conf := range config.Cfg.Mysql {
		engine, err := createMysqlEngine(conf)
		if err != nil {
			return fmt.Errorf("load SqlEngine failed: dbname(%s),err(%v)", dbName, err)
		}
		// engine.Debug()
		db, err := engine.DB()
		if err != nil {
			return err
		}

		db.SetMaxOpenConns(conf.MaxConnNum)
		db.SetMaxIdleConns(conf.MaxIdleConnNum)
		db.SetConnMaxIdleTime(time.Duration(conf.MaxConnLifeTime) * time.Second)
		db.SetConnMaxLifetime(time.Duration(conf.MaxConnLifeTime) * time.Second)

		engineManager[dbName] = engine
	}

	log.Printf("engineManager=%v", engineManager)
	return nil
}

func createMysqlEngine(conf config.MySqlConfig) (*gorm.DB, error) {
	//here can use xorm.EngineGroup  for slave db.
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=true",
		conf.User, conf.Passwd, "tcp", conf.Host,
		conf.Port, conf.DBName, conf.ConnTimeout, conf.ReadTimeout, conf.WriteTimeout)

	read_dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=true",
		conf.User, conf.Passwd, "tcp", conf.ReadHost,
		conf.Port, conf.DBName, conf.ConnTimeout, conf.ReadTimeout, conf.WriteTimeout)

	engine, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.LogLevel(conf.LogLevel))})

	if err = engine.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dsn)},
		Replicas: []gorm.Dialector{mysql.Open(read_dsn)},
	})); err != nil {
		return nil, err
	}

	if err = engine.Use(&MysqlPrometheus{}); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return engine, nil
}

func GetMysqlEngine(db string, ctx context.Context) (*gorm.DB, error) {
	if engine, ok := engineManager[db]; ok && engine != nil {
		ctxPm := ctx
		if gctx, ok := ctx.(*gin.Context); ok {
			ctxPm = gctx.Request.Context()
		}
		engine.Statement.Context = ctxPm

		return engine, nil
	}
	return nil, ErrNotFoundEngine
}
