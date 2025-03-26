package main

import (
	"fmt"
	"github.com/KNICEX/InkFlow/ioc"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/worker"
	"net/http"
)

func main() {
	initViper()
	app := ioc.InitApp()
	initPrometheus()

	for _, consumer := range app.Consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}
	}

	for _, w := range app.Workers {
		go func() {
			err := w.Run(worker.InterruptCh())
			if err != nil {
				fmt.Println("worker run err", err)
			}
		}()
	}

	app.Server.Run(":8080")

}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func initViper() {

	// --config=xxx.yaml
	configFile := pflag.String("config", "", "specify config file")
	pflag.Parse()
	if configFile != nil && *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.WatchConfig()
}
