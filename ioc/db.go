package ioc

import (
	"fmt"
	"github.com/KNICEX/InkFlow/pkg/gormx"
	"github.com/KNICEX/InkFlow/pkg/logx"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"strings"
	"time"
)

func InitDB(l logx.Logger) *gorm.DB {
	type Config struct {
		DSN string `mapstructure:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("postgres", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.New(gormLogFunc(l.Debug), logger.Config{
			SlowThreshold:             time.Millisecond * 50,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Warn,
		}),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "ink_flow",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.Postgres{
				VariableNames: []string{"Threads_running"},
			},
		},
	}))
	err = db.Use(tracing.NewPlugin(
		tracing.WithDBName("ink_flow"),
		tracing.WithoutMetrics()))
	if err != nil {
		panic(err)
	}

	vector := prometheus2.NewSummaryVec(prometheus2.SummaryOpts{
		Namespace: "ink_flow",
		Subsystem: "db",
		Name:      "gorm_sql_time",
		Help:      "gorm sql time",
		ConstLabels: map[string]string{
			"db": "ink_flow",
		},
		Objectives: map[float64]float64{
			0.5:   0.05,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})
	prometheus2.MustRegister(vector)
	cbs := gormx.NewSqlTimeCallbacks(vector, func(typ string, table string, db *gorm.DB) []string {
		return []string{typ, table}
	})
	if err = db.Use(cbs); err != nil {
		panic(err)
	}

	return db
}

type gormLogFunc func(msg string, fields ...logx.Field)

func (g gormLogFunc) Printf(format string, args ...any) {
	sql := fmt.Sprintf(format, args...)
	if strings.Contains(sql, "pg_") || strings.Contains(sql, "information_schema") {
		return
	}
	g("GORM: ", logx.Any("log", sql))
}
