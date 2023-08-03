package drivers

import (
	"prometheus-test/metrics"
	"time"

	"gorm.io/gorm"
)

type MysqlPrometheus struct {
}

func (m *MysqlPrometheus) Name() string {
	return "gorm:MysqlPrometheus"
}

func (m *MysqlPrometheus) Initialize(db *gorm.DB) (err error) {
	// before database operation
	db.Callback().Create().Before("gorm:create").Register("gorm_create", m.BeforeCallback("create"))
	db.Callback().Query().Before("gorm:query").Register("gorm_create", m.BeforeCallback("query"))
	db.Callback().Update().Before("gorm:update").Register("gorm_create", m.BeforeCallback("update"))
	db.Callback().Delete().Before("gorm:delete").Register("gorm_create", m.BeforeCallback("delete"))
	db.Callback().Row().Before("gorm:row").Register("gorm_create", m.BeforeCallback("row"))
	db.Callback().Raw().Before("gorm:raw").Register("gorm_create", m.BeforeCallback("raw"))

	// after database operation
	db.Callback().Create().After("gorm:create").Register("gorm_end", m.AfterCallback())
	db.Callback().Query().After("gorm:query").Register("gorm_end", m.AfterCallback())
	db.Callback().Update().After("gorm:update").Register("gorm_end", m.AfterCallback())
	db.Callback().Delete().After("gorm:delete").Register("gorm_end", m.AfterCallback())
	db.Callback().Row().After("gorm:row").Register("gorm_end", m.AfterCallback())
	db.Callback().Raw().After("gorm:raw").Register("gorm_end", m.AfterCallback())

	return
}
func (m *MysqlPrometheus) BeforeCallback(operation string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		db.Set("gormBegin", time.Now().UnixMilli())
		db.Set("gormCMD", operation)
	}
}
func (m *MysqlPrometheus) AfterCallback() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		cmd, cmdOK := db.Get("gormCMD")
		if !cmdOK {
			cmd = "all"
		}
		begin, bOK := db.Get("gormBegin")
		if !bOK {
			return
		}

		err := db.Error

		timeCost := time.Now().UnixMilli() - begin.(int64)
		metrics.UpdateDB(db.Statement.Table, cmd.(string), timeCost, err)
		metrics.UpdateDBQPS(db.Statement.Table, cmd.(string), err, 1)

	}
}
