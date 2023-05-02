package migrations

import (
	"github.com/bachelor/internal/config"
	"github.com/bachelor/internal/model"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"time"
)

func Run() error {
	var (
		vp  *viper.Viper
		err error
	)

	if vp, err = config.InitConfig(filepath.Join("configs"), "config"); err != nil {
		log.Fatal(err)
	}

	conn := pg.Connect(&pg.Options{
		User:     vp.GetString("postgres.login"),
		Password: vp.GetString("postgres.password"),
		Addr:     vp.GetString("postgres.host") + ":" + vp.GetString("postgres.port"),
		Database: vp.GetString("postgres.name"),
	})

	if err = conn.Model(&model.FiltrationRule{}).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}
	if err = conn.Model(&model.TransformationRule{}).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}
	if err = conn.Model(&model.DeduplicationRule{}).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}
	if err = conn.Model(&model.EnrichmentRule{}).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}
	if err = conn.Model(&model.ActionRule{}).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}

	filtrationRule := &model.FiltrationRule{
		Field:    "filtration_field",
		Function: "equals",
		Value:    "filtration_value",
	}

	transformationRule := &model.TransformationRule{
		Field:       "transformation_field",
		TargetField: "transformation_target_field",
	}

	deduplicationRule := &model.DeduplicationRule{
		Field: "",
		Value: map[string]struct{}{"deduplication_value": {}},
		Start: time.Now(),
		End:   time.Now().Add(time.Hour * 24),
	}

	enrichmentRule := &model.EnrichmentRule{
		Field: "enrichment_field",
		Value: "enrichment_value",
	}

	actionRule := &model.ActionRule{
		BootstrapServers: "localhost:29092",
		Key:              "action_field",
	}

	if _, err = conn.Model(filtrationRule).Insert(); err != nil {
		return err
	}
	if _, err = conn.Model(transformationRule).Insert(); err != nil {
		return err
	}
	if _, err = conn.Model(deduplicationRule).Insert(); err != nil {
		return err
	}
	if _, err = conn.Model(enrichmentRule).Insert(); err != nil {
		return err
	}
	if _, err = conn.Model(actionRule).Insert(); err != nil {
		return err
	}

	return nil
}
