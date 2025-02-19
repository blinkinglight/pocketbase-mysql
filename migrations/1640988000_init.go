//go:build !mysql

// Package migrations contains the system PocketBase DB migrations.
package migrations

import (
	"fmt"

	"github.com/blinkinglight/pocketbase-mysql/daos"
	"github.com/blinkinglight/pocketbase-mysql/models"
	"github.com/blinkinglight/pocketbase-mysql/models/schema"
	"github.com/blinkinglight/pocketbase-mysql/tools/migrate"
	"github.com/pocketbase/dbx"
)

var AppMigrations migrate.MigrationsList

// Register is a short alias for `AppMigrations.Register()`
// that is usually used in external/user defined migrations.
func Register(
	up func(db dbx.Builder) error,
	down func(db dbx.Builder) error,
	optFilename ...string,
) {
	AppMigrations.Register(up, down, optFilename...)
}

func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		_, tablesErr := db.NewQuery(`
			CREATE TABLE _admins (
				[[id]]              TEXT PRIMARY KEY,
				[[avatar]]          INTEGER DEFAULT 0 NOT NULL,
				[[email]]           TEXT UNIQUE NOT NULL,
				[[tokenKey]]        TEXT UNIQUE NOT NULL,
				[[passwordHash]]    TEXT NOT NULL,
				[[lastResetSentAt]] TEXT DEFAULT '' NOT NULL,
				[[created]]         TEXT DEFAULT '' NOT NULL,
				[[updated]]         TEXT DEFAULT '' NOT NULL
			);

			CREATE TABLE {{_users}} (
				[[id]]                     TEXT PRIMARY KEY,
				[[verified]]               BOOLEAN DEFAULT FALSE NOT NULL,
				[[email]]                  TEXT UNIQUE NOT NULL,
				[[tokenKey]]               TEXT UNIQUE NOT NULL,
				[[passwordHash]]           TEXT NOT NULL,
				[[lastResetSentAt]]        TEXT DEFAULT '' NOT NULL,
				[[lastVerificationSentAt]] TEXT DEFAULT '' NOT NULL,
				[[created]]                TEXT DEFAULT '' NOT NULL,
				[[updated]]                TEXT DEFAULT '' NOT NULL
			);

			CREATE TABLE {{_collections}} (
				[[id]]         TEXT PRIMARY KEY,
				[[system]]     BOOLEAN DEFAULT FALSE NOT NULL,
				[[name]]       TEXT UNIQUE NOT NULL,
				[[schema]]     JSON DEFAULT '[]' NOT NULL,
				[[listRule]]   TEXT DEFAULT NULL,
				[[viewRule]]   TEXT DEFAULT NULL,
				[[createRule]] TEXT DEFAULT NULL,
				[[updateRule]] TEXT DEFAULT NULL,
				[[deleteRule]] TEXT DEFAULT NULL,
				[[created]]    TEXT DEFAULT '' NOT NULL,
				[[updated]]    TEXT DEFAULT '' NOT NULL
			);

			CREATE TABLE {{_params}} (
				[[id]]      TEXT PRIMARY KEY,
				[[key]]     TEXT UNIQUE NOT NULL,
				[[value]]   JSON DEFAULT NULL,
				[[created]] TEXT DEFAULT '' NOT NULL,
				[[updated]] TEXT DEFAULT '' NOT NULL
			);
		`).Execute()
		if tablesErr != nil {
			return tablesErr
		}

		// inserts the system profiles collection
		// -----------------------------------------------------------
		profileOwnerRule := fmt.Sprintf("%s = @request.user.id", models.ProfileCollectionUserFieldName)
		collection := &models.Collection{
			Name:       models.ProfileCollectionName,
			System:     true,
			CreateRule: &profileOwnerRule,
			ListRule:   &profileOwnerRule,
			ViewRule:   &profileOwnerRule,
			UpdateRule: &profileOwnerRule,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     models.ProfileCollectionUserFieldName,
					Type:     schema.FieldTypeUser,
					Unique:   true,
					Required: true,
					System:   true,
					Options: &schema.UserOptions{
						MaxSelect:     1,
						CascadeDelete: true,
					},
				},
				&schema.SchemaField{
					Name:    "name",
					Type:    schema.FieldTypeText,
					Options: &schema.TextOptions{},
				},
				&schema.SchemaField{
					Name: "avatar",
					Type: schema.FieldTypeFile,
					Options: &schema.FileOptions{
						MaxSelect: 1,
						MaxSize:   5242880,
						MimeTypes: []string{
							"image/jpg",
							"image/jpeg",
							"image/png",
							"image/svg+xml",
							"image/gif",
						},
					},
				},
			),
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		tables := []string{
			"_params",
			"_collections",
			"_users",
			"_admins",
			models.ProfileCollectionName,
		}

		for _, name := range tables {
			if _, err := db.DropTable(name).Execute(); err != nil {
				return err
			}
		}

		return nil
	})
}
