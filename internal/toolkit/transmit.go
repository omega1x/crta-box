package toolkit

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Make actual data transfer from *CRTA-BOX* server to target *ClickHouse*-database
func Transmit(
	box_host string, box_port uint16, box_database string, box_username string, box_password string,
	house_host string, house_port uint16, house_database string, house_table string, house_username string, house_password string,
) (
	house_max_timestamp int64,
	stats_count_row uint,
	stats_count_row_fail_Scan uint,
	stats_count_row_fail_OutOfRange uint,
	stats_count_row_fail_Batch uint,
	stats_timestamp_first int64,
	stats_timestamp_last int64,
	err error,
) {

	// Scanned tags: see specification in *.share/ch/create-table.sql*
	var (
		timestamp int64 // use int64 to be compatible with [time.Unix()](https://pkg.go.dev/time#Unix)

		id_unit   uint64
		id_box    uint16
		id_spot   uint8
		id_sensor uint32

		sensor_serial   uint32
		sensor_revision uint8

		value_status    int8
		value_lf        uint16
		value_hf        uint16
		value_signal    uint16
		value_ratio     int8
		preset_signal   uint8
		preset_ratio    int8
		threshold_halt  uint8
		threshold_fault uint16

		sensor_voltage uint8
		relay_delay    uint8
		filter_gain    uint8

		StatusCode uint32 // absent in *CRTA-BOX*
	)

	house_db, house_ctx, err := driverHouse(house_host, house_port, house_database, house_username, house_password)
	if err != nil {
		err = fmt.Errorf("fail perform *Clickhouse* open operation " + err.Error())
		return
	}
	defer house_db.Close()
	defer house_ctx.Done()
	house_row := house_db.QueryRow(house_ctx, fmt.Sprintf(SQL_TEMPLATE_MAX_TIMESTAMP, house_table))
	var house_max_time time.Time
	err = house_row.Scan(&house_max_time)
	if err != nil {
		err = fmt.Errorf("fail perform row scan operation for values from *Clickhouse* " + err.Error())
		return
	}

	house_max_timestamp = house_max_time.Unix()
	if house_max_timestamp < GENESIS_TIMESTAMP {
		// The target table is highly likely empty
		house_max_timestamp = time.Now().Add(time.Minute * DEFAULT_HOUSE_TIME_LAG).Unix()
	}

	box_db, err := sql.Open(BOX_DATABASE_SYSTEM, fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", box_host, box_port, box_database, box_username, box_password))
	if err != nil {
		err = fmt.Errorf("fail perform *CRTA-BOX* open operation " + err.Error())
		return
	}
	defer box_db.Close()

	box_rows, err := box_db.Query(fmt.Sprintf(`
	SELECT 
	   timestamp
	  ,id_unit
	  ,id_box
	  ,id_spot
	  ,id_sensor
	  ,sensor_serial
	  ,sensor_revision
	  ,value_status
	  ,value_lf
	  ,value_hf
	  ,value_signal
	  ,value_ratio
	  ,abs(preset_signal) AS preset_signal  -- adapt always negative values of initial signal to uint8-type
	  ,preset_ratio
	  ,threshold_halt
	  ,threshold_fault
	  ,sensor_voltage
	  ,relay_delay
	  ,filter_gain
	FROM box_view
		WHERE 
				timestamp > %d 
			AND timestamp < (
				SELECT MAX(timestamp) - 1 FROM box_view WHERE timestamp >= %d
			)
	ORDER BY timestamp
	`, house_max_timestamp, house_max_timestamp))

	if err != nil {
		err = fmt.Errorf("fail perform *CRTA-BOX* SQL-query operation " + err.Error())
		return
	}
	defer box_rows.Close()

	house_batch, err := house_db.PrepareBatch(
		house_ctx, fmt.Sprintf("INSERT INTO %s.%s", house_database, house_table),
	)

	if err != nil {
		err = fmt.Errorf("fail prepare batch in *ClickHouse* " + err.Error())
		return
	}

	for box_rows.Next() {
		err = box_rows.Scan(
			&timestamp,
			&id_unit,
			&id_box,
			&id_spot,
			&id_sensor,

			&sensor_serial,
			&sensor_revision,

			&value_status,
			&value_lf,
			&value_hf,
			&value_signal,
			&value_ratio,
			&preset_signal,
			&preset_ratio,
			&threshold_halt,
			&threshold_fault,

			&sensor_voltage,
			&relay_delay,
			&filter_gain,
		)
		if err != nil {
			err = fmt.Errorf("fail perform scan operation for values form *CRTA-BOX* " + err.Error())
			stats_count_row_fail_Scan++
			continue
		}

		// Signal diagnostics and modifications in current row
		StatusCode = OPC_CODE_GoodEntryInserted

		if value_status < PRESET_MIN_VALUE_STATUS || value_status > PRESET_MAX_VALUE_STATUS {
			value_status = PRESET_INCORRECT_VALUE_STATUS
			StatusCode = OPC_CODE_BadOutOfRange
		}

		/*
			if value_ratio < 0 || value_grade < 0 || preset_ratio < 0 || preset_grade < 0 || value_noise < 0 ||
					preset_noise < 0 || value_gain < 0 || value_high < 0 || value_low < 0 {
					StatusCode = conf.OPC_CODE_BadOutOfRange
				}
		*/

		// Statistics calculations
		if stats_count_row == 0 {
			stats_timestamp_first = timestamp
		}
		/*if StatusCode == conf.OPC_CODE_BadOutOfRange {
			stats_count_row_fail_OutOfRange++
		}
		*/

		// Output preparations:
		err = house_batch.Append(
			uuid.New(),
			time.Unix(timestamp, 0),

			id_unit,
			id_box,
			id_spot,
			id_sensor,

			sensor_serial,
			sensor_revision,

			value_status,
			value_lf,
			value_hf,
			value_signal,
			value_ratio,

			preset_signal,
			preset_ratio,

			threshold_halt,
			threshold_fault,

			sensor_voltage,
			relay_delay,
			filter_gain,

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
