package selection

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jukeizu/selection/selection/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "selection"
)

type Repository interface {
	Migrate() error
	CreateSelection(Selection) error
	Selection(appId, userId, serverId string) (Selection, error)
}

type repository struct {
	Db *sql.DB
}

func NewRepository(url string) (Repository, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	r := repository{
		Db: db,
	}

	return &r, nil
}

func (r *repository) Migrate() error {
	_, err := r.Db.Exec(`CREATE DATABASE IF NOT EXISTS ` + DatabaseName)
	if err != nil {
		return err
	}

	g, err := gossage.New(r.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(
		migrations.CreateTableSelection20190415004138{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) CreateSelection(selection Selection) error {
	q := `INSERT INTO selection (appId, userId, serverId, options)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (appId, userId, serverId)
		DO UPDATE SET options = excluded.options, updated = now()`

	options, err := json.Marshal(selection.Options)
	if err != nil {
		return fmt.Errorf("Could not marshal options to JSON: %s", err)
	}

	_, err = r.Db.Exec(q, selection.AppId, selection.UserId, selection.ServerId, options)

	return err
}

func (r *repository) Selection(appId, userId, serverId string) (Selection, error) {
	//	q := `SELECT ()`

	return Selection{}, nil

}
