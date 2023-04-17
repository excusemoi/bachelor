package db

import (
	"errors"
	"github.com/bachelor/internal/model/filter"
	"time"
)

func (db *Db) GetLatestFiltrationRule() (*filter.FiltrationRule, error) {
	var fr filter.FiltrationRule
	err := db.client.Model(&fr).Order("updated_at DESC").Last()
	return &fr, err
}

func (db *Db) GetFiltrationRule(id int) (*filter.FiltrationRule, error) {
	var fr filter.FiltrationRule
	err := db.client.Model(&fr).Where("id = ?", id).Select()
	return &fr, err
}

func (db *Db) CreateFiltrationRule(fr *filter.FiltrationRule) error {
	var err error
	if fr == nil {
		return errors.New("filtrationRule must not be nil")
	}
	fr.Id = 0
	if err = db.CreateFilter(fr.Filter); err != nil {
		return err
	}
	fr.FilterId = fr.Filter.Id
	_, err = db.client.Model(fr).Insert()
	return err
}

func (db *Db) UpdateFiltrationRule(fr *filter.FiltrationRule) error {
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

func (db *Db) CreateFilter(f *filter.Filter) error {
	if f == nil {
		return errors.New("filter must not be nil")
	}
	f.Id = 0
	_, err := db.client.Model(f).Insert()
	return err
}
