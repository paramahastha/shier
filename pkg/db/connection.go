package db

import "github.com/go-pg/pg"

var dbConn *pg.DB

func NewConnection(c *Config) (*pg.DB, error) {
	opt, err := pg.ParseURL(c.URL)
	if err != nil {
		return nil, err
	}

	dbConn = pg.Connect(opt)

	// ensure database connection is successfully
	_, err = dbConn.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func GetConnection() *pg.DB {
	return dbConn
}
