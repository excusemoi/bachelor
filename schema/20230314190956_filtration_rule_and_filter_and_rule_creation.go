package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		sql := `
			create table if not exists filter (
                                      id serial primary key,
                                      filter text
			);
		`
		if _, err := db.Exec(sql); err != nil {
			return err
		}

		sql = `
			create table if not exists rule (
                                    id serial primary key,
                                    rule text
			);
		`
		if _, err := db.Exec(sql); err != nil {
			return err
		}

		sql = `
			create table if not exists filtration_rule (
				id serial primary key,
				filter_id integer constraint filter_id_constraint references filter on delete cascade,
				rule_id integer constraint rule_id_constraint references rule on delete cascade,
				filter_field text,
				filter_function text,
				filter_value text
			    updated_at Timestamp not null
			);
		`
		if _, err := db.Exec(sql); err != nil {
			return err
		}

		return nil
	}

	down := func(db orm.DB) error {
		sql := `
			drop table if exists filtration_rule;
		`
		if _, err := db.Exec(sql); err != nil {
			return err
		}

		sql = `
			drop table if exists filter;
		`
		if _, err := db.Exec(sql); err != nil {
			return err
		}

		sql = `
			drop table if exists rule;
		`
		if _, err := db.Exec(sql); err != nil {
			return err
		}
		return nil
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20230314190956_filtration_rule_and_filter_and_rule_creation", up, down, opts)
}
