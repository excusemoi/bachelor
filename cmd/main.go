package main

import (
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/components/actor"
	"github.com/bachelor/internal/components/deduplicator"
	"github.com/bachelor/internal/components/enricher"
	"github.com/bachelor/internal/components/filter"
	"github.com/bachelor/internal/components/source"
	"github.com/bachelor/internal/components/transformer"
	"github.com/bachelor/internal/config"
	"github.com/bachelor/internal/migrations"
	"github.com/bachelor/internal/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	var (
		ac         = &components.AbstractComponent[model.AbstractRule]{Wg: &sync.WaitGroup{}}
		src        = &source.Source{}
		fl         = &filter.Filter{}
		tf         = &transformer.Transformer{}
		dd         = &deduplicator.Deduplicator{}
		enr        = &enricher.Enricher{}
		act        = &actor.Actor{}
		configName = "config"
		cfg        = &viper.Viper{}
		err        error
	)

	if cfg, err = config.InitConfig("configs", configName); err != nil {
		log.Println("config")
		log.Fatal(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err = migrations.Run(); err != nil {
		log.Println("migrations")
		log.Fatal(err)
	}

	if err = src.Init(filepath.Join("internal", "components", "source", "configs"), configName); err != nil {
		log.Println("source")
		log.Fatal(err)
	}
	if err = tf.Init(filepath.Join("internal", "components", "transformer", "configs"), configName); err != nil {
		log.Println("transformer")
		log.Fatal(err)
	}
	if err = fl.Init(filepath.Join("internal", "components", "filter", "configs"), configName); err != nil {
		log.Println("filter")
		log.Fatal(err)
	}
	if err = dd.Init(filepath.Join("internal", "components", "deduplicator", "configs"), configName); err != nil {
		log.Println("deduplicator")
		log.Fatal(err)
	}
	if err = enr.Init(filepath.Join("internal", "components", "enricher", "configs"), configName); err != nil {
		log.Println("enricher")
		log.Fatal(err)
	}
	if err = act.Init(filepath.Join("internal", "components", "actor", "configs"), configName, func(s string) string { return s }); err != nil {
		log.Println("actor")
		log.Fatal(err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(cfg.GetString("prometheus.host")+":"+cfg.GetString("prometheus.port"),
			nil); err != nil {
			log.Println("prometheus")
			log.Fatal(err)
		}
	}()

	ac.SetNext(src).SetNext(fl).SetNext(tf).SetNext(dd).SetNext(enr).SetNext(act)
	ac.RunPipeline()
}
