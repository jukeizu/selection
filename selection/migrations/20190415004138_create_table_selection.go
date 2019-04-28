package migrations

import (
	"database/sql"
)

type CreateTableSelection20190415004138 struct{}

func (m CreateTableSelection20190415004138) Version() string {
	return "20190415004138_CreateTableSelection"
}

func (m CreateTableSelection20190415004138) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS selection (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			appId STRING NOT NULL DEFAULT '',
			instanceId STRING NOT NULL DEFAULT '',
			userId STRING NOT NULL DEFAULT '',
			serverId STRING NOT NULL DEFAULT '',
			options JSONB NOT NULL,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			UNIQUE (appId, userId, serverId)
		)`)

	return err
}

func (m CreateTableSelection20190415004138) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE selection`)
	return err
}
