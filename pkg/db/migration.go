package db

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/paramahastha/shier/assets"
	migrate "github.com/rubenv/sql-migrate"
)

func Migrate(dbURL string) error {
	parsedDbUrl, err := url.Parse(dbURL)
	if err != nil {
		return fmt.Errorf("error while parsing URL %s", err.Error())
	}

	db, err := gorm.Open(parsedDbUrl.Scheme, dbURL)
	if err != nil {
		return fmt.Errorf("failed connect to database %s", err.Error())
	}
	defer db.Close()

	migrations := &migrate.AssetMigrationSource{
		Asset:    assets.Asset,
		AssetDir: assets.AssetDir,
		Dir:      "assets/sql",
	}

	migrate.SetTable("migrations")

	_, err = migrate.Exec(
		db.DB(),
		parsedDbUrl.Scheme,
		migrations,
		migrate.Up,
	)

	return err
}
