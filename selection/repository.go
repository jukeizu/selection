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
	Selection(appId, instanceId, userId, serverId string) (Selection, error)
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
	q := `INSERT INTO selection (appId, instanceId, userId, serverId, options)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (appId, userId, serverId)
		DO UPDATE SET instanceId = excluded.instanceId, options = excluded.options, updated = now()`

	options, err := json.Marshal(selection.Batches)
	if err != nil {
		return fmt.Errorf("could not marshal options to JSON: %s", err)
	}

	_, err = r.Db.Exec(q, selection.AppId, selection.InstanceId, selection.UserId, selection.ServerId, options)

	return err
}

func (r *repository) Selection(appId, instanceId, userId, serverId string) (Selection, error) {
	q := `SELECT appId, instanceId, userId, serverId, options FROM selection
	WHERE appId = $1 AND instanceId = $2 AND userId = $3 AND serverId = $4`

	selection := Selection{}

	jsonOptions := []byte{}

	err := r.Db.QueryRow(q, appId, instanceId, userId, serverId).Scan(
		&selection.AppId,
		&selection.InstanceId,
		&selection.UserId,
		&selection.ServerId,
		&jsonOptions,
	)
	if err != nil {
		return Selection{}, err
	}

	err = json.Unmarshal(jsonOptions, &selection.Batches)
	if err != nil {
		return Selection{}, fmt.Errorf("could not unmarshal JSON to options: %s", err)
	}

	return selection, nil
}
