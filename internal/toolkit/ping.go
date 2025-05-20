package toolkit

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Check connection with *CRTA-BOX* server
func PingBox(host string, port uint16, database string, username string, password string) (string, error) {
	var box_max_timestamp int64
	var note_info string

	db, err := sql.Open(BOX_DATABASE_SYSTEM, fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, database, username, password))
	if err != nil {
		return fmt.Sprintf("fail connect to *CRTA-BOX* server (%s) ", BOX_DATABASE_SYSTEM), err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Sprintf("fail pinging *CRTA-BOX* database server (%s): server might be unavailable", BOX_DATABASE_SYSTEM), err
	}

	row := db.QueryRow(fmt.Sprintf(SQL_TEMPLATE_MAX_TIMESTAMP, BOX_TABLE_ORIGIN))
	err = row.Err()

	if err != nil {
		return fmt.Sprintf("pinged but fail query table [%s.%s]: table might be unavailable", database, BOX_TABLE_ORIGIN), err
	}

	err = row.Scan(&box_max_timestamp)
	if err != nil {
		return fmt.Sprintf("pinged but fail get maximum timestamp from table [%s.%s]: column `timestamp` might have wrong type", database, BOX_TABLE_ORIGIN), err
	}

	if box_max_timestamp < time.Now().Unix()-BOX_DATA_EXPIRY_PERIOD {
		note_info = " data in *CRTA-BOX* is may be too old or absent "
	}

	return fmt.Sprintf("Maximum timestamp is %d (%s)"+note_info, box_max_timestamp, time.Unix(box_max_timestamp, 0).Format(time.RFC1123)), nil
}

// Check connection with the *ClickHouse*-database
func PingHouse(host string, port uint16, database string, table string, username string, password string) (string, error) {
	const L_EMPTY_NAME = ""
	var (
		house_max_timestamp time.Time
		note_info           string
	)

	db, ctx, err := driverHouse(host, port, database, username, password)
	if err != nil {
		return "fail connect to *ClickHouse*-database ", err
	}
	defer db.Close()
	defer ctx.Done()

	err = db.Ping(ctx)
	if err != nil {
		return "fail pinging *ClickHouse*-database. Server might be unavailable ", err
	}

	if table == L_EMPTY_NAME {
		return ". But data availability cannot be confirmed since target table is not specified", nil
	}

	row := db.QueryRow(ctx, fmt.Sprintf("SELECT 1 FROM %s.%s", database, table))
	err = row.Err()

	if err != nil {
		return fmt.Sprintf("pinged but fail query table [%s.%s]: table might be unavailable ", database, table), err
	}

	row = db.QueryRow(ctx, fmt.Sprintf(SQL_TEMPLATE_MAX_TIMESTAMP, table))
	err = row.Scan(&house_max_timestamp)
	if err != nil {
		return fmt.Sprintf("pinged but fail get maximum timestamp from table [%s.%s]: column `timestamp` might the have wrong type ", database, table), err
	}

	if house_max_timestamp.Unix() < GENESIS_TIMESTAMP {
		note_info = fmt.Sprintf(". Table [%s.%s] might be empty", database, table)
	}

	return fmt.Sprintf("Maximum timestamp is %d (%s)"+note_info, house_max_timestamp.Unix(), house_max_timestamp.Format(time.RFC1123)), nil
}
