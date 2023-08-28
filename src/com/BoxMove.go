package com

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/omega1x/crta-box/src/conf"
)

// Make actual data transfer from *CRTA-BOX* to target ClickHouse
func BoxMove(
	box_host string, box_port uint16, box_database string, box_username string, box_password string,
	house_host string, house_port uint16, house_database string, house_username string, house_password string,
) (
	house_max_timestamp uint64,
	stats_count_row uint,
	stats_count_row_fail_Scan uint,
	stats_count_row_fail_OutOfRange uint,
	stats_count_row_fail_Batch uint,
	stats_timestamp_first int64,
	stats_timestamp_last int64,
	err error,
) {

	var (

		// Scanned tags: see specification in *./share/ch__create-table__log_box3.sql*
		timestamp    int64
		id_unit      uint64
		id_box       uint16
		id_spot      uint8
		id_sensor    uint32
		value_status int16
		value_ratio  int16
		value_grade  int16
		preset_ratio int16
		preset_grade int16
		value_noise  int16
		preset_noise int16
		value_gain   int16
		value_high   int16
		value_low    int16
		StatusCode   uint32 // absent in *CRTA-BOX*
	)
	house_db, err := HouseDrive(house_host, house_port, house_database, house_username, house_password)
	if err != nil {
		return
	}

	house_ctx := context.Background()
	defer house_ctx.Done()

	house_row := house_db.QueryRow(house_ctx, conf.SQL_SELECT_HOUSE_MAX_TIMESTAMP)
	defer house_db.Close()

	err = house_row.Scan(&house_max_timestamp)
	if err != nil {
		return
	}

	if house_max_timestamp < conf.GENESIS_TIMESTAMP {
		err = fmt.Errorf("there must be no actual data in *ClickHouse*")
		return
	}

	box_db, err := sql.Open(
		conf.BOX_DATABASE_SYSTEM,
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", box_host, box_port, box_database, box_username, box_password),
	)
	if err != nil {
		return
	}
	defer box_db.Close()

	box_rows, err := box_db.Query(fmt.Sprintf(conf.SQL_SELECT_BOX_DATA, house_max_timestamp, house_max_timestamp))
	if err != nil {
		return
	}
	defer box_rows.Close()

	house_batch, err := house_db.PrepareBatch(house_ctx, conf.SQL_INSERT_HOUSE_DATA)
	if err != nil {
		return
	}

	for box_rows.Next() {
		err = box_rows.Scan(
			&timestamp,
			&id_unit,
			&id_box,
			&id_spot,
			&id_sensor,
			&value_status,
			&value_ratio,
			&value_grade,
			&preset_ratio,
			&preset_grade,
			&value_noise,
			&preset_noise,
			&value_gain,
			&value_high,
			&value_low,
		)

		if err != nil {
			stats_count_row_fail_Scan++
			continue
		}

		// Signal diagnostics and modifications in current row
		StatusCode = conf.OPC_CODE_GoodEntryInserted
		if value_status < conf.TAG_PRESET_STATUS_MIN || value_status > conf.TAG_PRESET_STATUS_MAX {
			value_status = conf.TAG_PRESET_STATUS_INCORRECT
			StatusCode = conf.OPC_CODE_BadOutOfRange
		}
		if value_ratio < 0 || value_grade < 0 || preset_ratio < 0 || preset_grade < 0 || value_noise < 0 ||
			preset_noise < 0 || value_gain < 0 || value_high < 0 || value_low < 0 {
			StatusCode = conf.OPC_CODE_BadOutOfRange
		}

		// Statistics calculations
		if stats_count_row == 0 {
			stats_timestamp_first = timestamp
		}
		if StatusCode == conf.OPC_CODE_BadOutOfRange {
			stats_count_row_fail_OutOfRange++
		}

		// Output preparations:
		err = house_batch.Append(
			uuid.New(),

			time.Unix(timestamp, 0),

			id_unit,
			id_box,
			id_spot,

			id_sensor,

			int8(value_status),
			value_ratio,
			value_grade,
			preset_ratio,
			preset_grade,
			value_noise,
			preset_noise,
			value_gain,
			value_high,
			value_low,

			StatusCode,
		)
		if err != nil {
			stats_count_row_fail_Batch++
			continue
		}
		// So, a row has been processed:
		stats_count_row++
	}
	stats_timestamp_last = timestamp
	err = house_batch.Send()
	return
}
