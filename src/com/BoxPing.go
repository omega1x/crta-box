package com

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/omega1x/crta-box/src/conf"
)

// Check connection with *CRTA-BOX* instance with embedded PostgreSQL database
func BoxPing(host string, port uint16, database string, username string, password string) error {
	db, err := sql.Open(conf.BOX_DATABASE_SYSTEM, fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, database, username, password))
	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}
