package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type SqlTimeCallback struct {
	vector *prometheus.SummaryVec
	labels func(typ string, table string, db *gorm.DB) []string
}

func (c *SqlTimeCallback) Name() string {
	return "prometheus-sql-time-summary"
}

func (c *SqlTimeCallback) Initialize(db *gorm.DB) error {
	return c.RegisterAll(db)
}

func NewSqlTimeCallbacks(vector *prometheus.SummaryVec, labelsFunc func(typ string, table string, db *gorm.DB) []string) *SqlTimeCallback {
	res := &SqlTimeCallback{
		vector: vector,
		labels: func(typ string, table string, db *gorm.DB) []string {
			return []string{typ, table}
		},
	}
	if labelsFunc != nil {
		res.labels = labelsFunc
	}
	return res
}

func (c *SqlTimeCallback) RegisterAll(db *gorm.DB) error {
	err := db.Callback().Create().Before("*").
		Register("prometheus_before_create", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_after_create", c.After("create"))
	if err != nil {
		return err
	}
	err = db.Callback().Update().Before("*").
		Register("prometheus_before_update", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Update().After("*").
		Register("prometheus_after_update", c.After("update"))
	if err != nil {
		return err
	}
	err = db.Callback().Query().Before("*").
		Register("prometheus_before_query", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Query().After("*").
		Register("prometheus_after_query", c.After("query"))
	if err != nil {
		return err
	}
	err = db.Callback().Delete().Before("*").
		Register("prometheus_before_delete", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Delete().After("*").
		Register("prometheus_after_delete", c.After("delete"))
	if err != nil {
		return err
	}
	err = db.Callback().Row().Before("*").
		Register("prometheus_before_row", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Row().After("*").
		Register("prometheus_after_row", c.After("row"))
	if err != nil {
		return err
	}
	return nil
}

func (c *SqlTimeCallback) Before() func(*gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}

func (c *SqlTimeCallback) After(typ string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		c.vector.WithLabelValues(c.labels(typ, table, db)...).Observe(float64(time.Since(startTime).Milliseconds()))
	}
}
