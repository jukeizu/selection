package migrations

import (
	"database/sql"
)

type CreateTableSelectionOption20190415004504 struct{}

func (m CreateTableSelectionOption20190415004504) Version() string {
	return "20190415004504_CreateTableSelectionOption"
}

func (m CreateTableSelectionOption20190415004504) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS selection_option (
			id UUID NOT NULL DEFAULT gen_random_uuid(),
			selectionId UUID NOT NULL,
			selectionOptionIndex INT NOT NULL,
			optionId STRING NOT NULL DEFAULT '',
			content STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			PRIMARY KEY (selectionId, id),
			FOREIGN KEY (selectionId) REFERENCES selection (id) ON DELETE CASCADE
		) INTERLEAVE IN PARENT selection (selectionId)`)

	return err
}

func (m CreateTableSelectionOption20190415004504) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE selection_option`)
	return err
}
