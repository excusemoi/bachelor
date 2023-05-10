package db

import (
	"errors"
	"time"
)

func (db *Db[T]) Observe(m *T) {
	select {
	case <-time.After(10 * time.Second):
		{
			db.mx.Lock()
			m, _ = db.GetLatest(m)
			db.mx.Unlock()
		}
	}
}

func (db *Db[T]) GetLatest(model *T) (*T, error) {
	if err := db.client.Model(model).Order("updated_at DESC").Last(); err != nil {
		return nil, err
	}
	return model, nil
}

func (db *Db[T]) GetAll(model *[]*T) (*[]*T, error) {
	if model != nil && len(*model) != 0 {
		if err := db.client.Model((*model)[0]).Select(model); err != nil {
			return nil, err
		}
	}
	return model, nil
}

func (db *Db[T]) GetByID(id int) (*T, error) {
	var fr T
	err := db.client.Model(&fr).Where("id = ?", id).Select()
	return &fr, err
}

func (db *Db[T]) Create(m *T) error {
	var err error
	if m == nil {
		return errors.New("model must not be nil")
	}
	_, err = db.client.Model(m).Insert()
	return err
}
