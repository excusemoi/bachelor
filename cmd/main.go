package main

import (
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/components/actor"
	"github.com/bachelor/internal/components/deduplicator"
	"github.com/bachelor/internal/components/enricher"
	"github.com/bachelor/internal/components/filter"
	"github.com/bachelor/internal/components/source"
	"github.com/bachelor/internal/components/transformer"
	"github.com/bachelor/internal/migrations"
	"github.com/bachelor/internal/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

func main() {
	var err error

	if err = migrations.Run(); err != nil {
		log.Fatal(err)
	}

	ac := &components.AbstractComponent[model.AbstractRule]{Wg: &sync.WaitGroup{}}
	src := &source.Source{}
	fl := &filter.Filter{}
	tf := &transformer.Transformer{}
	dd := &deduplicator.Deduplicator{}
	enr := &enricher.Enricher{}
	act := &actor.Actor{}

	if err = src.Init(filepath.Join("internal", "components", "source", "configs")); err != nil {
		log.Fatal(err)
	}
	if err = tf.Init(filepath.Join("internal", "components", "transformer", "configs")); err != nil {
		log.Fatal(err)
	}
	if err = fl.Init(filepath.Join("internal", "components", "filter", "configs")); err != nil {
		log.Fatal(err)
	}
	if err = dd.Init(filepath.Join("internal", "components", "deduplicator", "configs")); err != nil {
		log.Fatal(err)
	}
	if err = enr.Init(filepath.Join("internal", "components", "enricher", "configs")); err != nil {
		log.Fatal(err)
	}
	if err = act.Init(filepath.Join("internal", "components", "actor", "configs"), func(s string) string { return s }); err != nil {
		log.Fatal(err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe("bachelor-app:8081", nil)
	}()

	ac.SetNext(src).SetNext(fl).SetNext(tf).SetNext(dd).SetNext(enr).SetNext(act)
	ac.RunPipeline()
}
