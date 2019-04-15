package migrations

import (
	"database/sql"
)

type CreateTableSelectionOptionMetadata20190415020132 struct{}

func (m CreateTableSelectionOptionMetadata20190415020132) Version() string {
	return "20190415020132_CreateTableSelectionOptionMetadata"
}

func (m CreateTableSelectionOptionMetadata20190415020132) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS selection_option_metadata (
			id UUID NOT NULL DEFAULT gen_random_uuid(),
			selectionId UUID NOT NULL,
			selectionOptionId UUID NOT NULL,
			key STRING NOT NULL DEFAULT '',
			value STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			PRIMARY KEY (selectionOptionId, selectionId, id),
			FOREIGN KEY (selectionId, selectionOptionId) REFERENCES selection_option (selectionId, id) ON DELETE CASCADE
		) INTERLEAVE IN PARENT selection_option (selectionOptionId, selectionId)`)

	return err
}

func (m CreateTableSelectionOptionMetadata20190415020132) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE selection_option_metadata`)
	return err
}
