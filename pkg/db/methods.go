package db

import (
	"errors"
	"github.com/bachelor/pkg/model"
	"time"
)

func (db *Db) GetLatestFiltrationRule() (*model.FiltrationRule, error) {
	var fr model.FiltrationRule
	err := db.client.Model(&fr).Order("updated_at DESC").Last()
	return &fr, err
}

func (db *Db) GetFiltrationRule(id int) (*model.FiltrationRule, error) {
	var fr model.FiltrationRule
	err := db.client.Model(&fr).Where("id = ?", id).Select()
	return &fr, err
}

func (db *Db) CreateFiltrationRule(fr *model.FiltrationRule) error {
	var err error
	if fr == nil {
		return errors.New("filtrationRule must not be nil")
	}
	fr.Id = 0
	if err = db.CreateFilter(fr.Filter); err != nil {
		return err
	}
	fr.FilterId = fr.Filter.Id
	if err = db.CreateRule(fr.Rule); err != nil {
		return err
	}
	fr.RuleId = fr.Rule.Id
	_, err = db.client.Model(fr).Insert()
	return err
}

func (db *Db) UpdateFiltrationRule(fr *model.FiltrationRule) error {
	if fr == nil {
		return errors.New("filtrationRule must not be nil")
	}
	if fr.Id == 0 {
		return errors.New("you must specify id")
	}
	fr.UpdatedAt = time.Now()
	_, err := db.client.Model(fr).Where("id = ?id", fr.Id).Update()
	return err
}

func (db *Db) CreateRule(r *model.Rule) error {
	if r == nil {
		return errors.New("rule must not be nil")
	}
	r.Id = 0
	_, err := db.client.Model(r).Insert()
	return err
}

func (db *Db) CreateFilter(f *model.Filter) error {
	if f == nil {
		return errors.New("filter must not be nil")
	}
	f.Id = 0
	_, err := db.client.Model(f).Insert()
	return err
}
