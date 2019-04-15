package selection

import (
	"database/sql"
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
		migrations.CreateTableSelection20190415004138{},
		migrations.CreateTableSelectionOption20190415004504{},
		migrations.CreateTableSelectionOptionMetadata20190415020132{},
	)
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) CreateSelection(selection Selection) error {
	q := `INSERT INTO selection (appId, userId, serverId)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.Db.QueryRow(q, selection.AppId, selection.UserId, selection.ServerId).Scan(
		&selection.Id,
	)
	if err != nil {
		return err
	}

	for _, option := range selection.Options {
		err := r.createSelectionOption(selection.Id, option)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) createSelectionOption(selectionId string, selectionOption SelectionOption) error {
	q := `INSERT INTO selection_option (
			selectionId, 
			selectionOptionIndex, 
			optionId, 
			content)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := r.Db.QueryRow(q,
		selectionId,
		selectionOption.SelectionOptionIndex,
		selectionOption.Option.OptionId,
		selectionOption.Option.Content,
	).Scan(&selectionOption.Id)

	if err != nil {
		return err
	}

	for k, v := range selectionOption.Option.Metadata {
		err := r.createSelectionOptionMetadata(selectionId, selectionOption.Id, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) createSelectionOptionMetadata(selectionId string, selectionOptionId string, key string, value string) error {
	q := `INSERT INTO selection_option_metadata (
			selectionId,
			selectionOptionId,
			key, 
			value)
		VALUES ($1, $2, $3, $4)`

	_, err := r.Db.Exec(q,
		selectionId,
		selectionOptionId,
		key,
		value,
	)

	return err
}
